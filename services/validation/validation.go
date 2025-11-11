package validation

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type RowAddr struct {
	ID     int64  `orm:"column(id)"`
	UUID   string `orm:"column(uuid)"`
	CityID int64  `orm:"column(city_id)"`
}

type Suggestion struct {
	Value             string          `json:"value"`
	UnrestrictedValue string          `json:"unrestricted_value"`
	Data              json.RawMessage `json:"data"`
}

type DadataClient interface {
	FindByIDWithRetry(id string, maxRetries int, backoffBase time.Duration) ([]Suggestion, error)
	MaxRetries() int
	BackoffBase() time.Duration
}

type PipelineOptions struct {
	RPS      int
	Workers  int
	LogEvery int
}

var DefaultPipelineOptions = PipelineOptions{
	RPS:      17,
	Workers:  12,
	LogEvery: 1000,
}

type PipelineOut struct {
	Responses map[string][]Suggestion
	Errors    map[string]error
	TotalSent int64
}

func RunPipeline(ctx context.Context, d DadataClient, ids []string, opts PipelineOptions) PipelineOut {
	rps := opts.RPS
	workers := opts.Workers
	logEvery := opts.LogEvery
	if rps <= 0 {
		rps = 1
	}
	if workers <= 0 {
		workers = 1
	}

	out := PipelineOut{
		Responses: make(map[string][]Suggestion, len(ids)),
		Errors:    make(map[string]error),
	}

	tokens := make(chan struct{}, rps)
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			select {
			case tokens <- struct{}{}:
			default:
			}
		}
	}()

	type job struct{ id string }
	type result struct {
		id   string
		sugs []Suggestion
		err  error
	}

	jobs := make(chan job, 1024)
	results := make(chan result, 1024)

	var wg sync.WaitGroup
	workerFn := func() {
		defer wg.Done()
		for j := range jobs {
			select {
			case <-ctx.Done():
				return
			case <-tokens:
			}
			atomic.AddInt64(&out.TotalSent, 1)
			sugs, err := d.FindByIDWithRetry(j.id, d.MaxRetries(), d.BackoffBase())
			if err != nil {
				select {
				case results <- result{id: j.id, err: err}:
				case <-ctx.Done():
				}
				continue
			}
			select {
			case results <- result{id: j.id, sugs: sugs}:
			case <-ctx.Done():
			}
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go workerFn()
	}

	go func() {
		defer close(jobs)
		for _, id := range ids {
			select {
			case <-ctx.Done():
				return
			case jobs <- job{id: id}:
			}
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		total := 0
		for total < len(ids) {
			select {
			case <-ctx.Done():
				return
			case r := <-results:
				if r.id == "" {
					continue
				}
				if r.err != nil {
					out.Errors[r.id] = r.err
				} else {
					out.Responses[r.id] = r.sugs
				}
				total++
				if logEvery > 0 && total%logEvery == 0 {
					log.Printf("[pipeline] got %d/%d (sent=%d, errors=%d, rps~%d, workers=%d)",
						total, len(ids), out.TotalSent, len(out.Errors), rps, workers)
				}
			}
		}
	}()
	<-done
	wg.Wait()
	return out
}

func LoadRowAddrsByIDs(o orm.Ormer, ids []int64) ([]RowAddr, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ph := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := range ids {
		ph[i] = "?"
		args[i] = ids[i]
	}
	q := "SELECT id, uuid, city_id FROM avileads_geo.geo_address WHERE id IN (" + strings.Join(ph, ",") + ")"
	var rows []RowAddr
	_, err := o.Raw(q, args...).QueryRows(&rows)
	return rows, err
}

func ExtractUUIDs(rows []RowAddr) []string {
	seen := make(map[string]struct{}, len(rows))
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		if _, ok := seen[r.UUID]; ok {
			continue
		}
		seen[r.UUID] = struct{}{}
		out = append(out, r.UUID)
	}
	return out
}

func GroupRowAddrsByCity(rows []RowAddr) map[int64][]RowAddr {
	m := make(map[int64][]RowAddr, 16)
	for _, r := range rows {
		m[r.CityID] = append(m[r.CityID], r)
	}
	return m
}

func FilterUnvalidatedRowAddrs(o orm.Ormer, rows []RowAddr) ([]RowAddr, []string, int, error) {
	if len(rows) == 0 {
		return nil, nil, 0, nil
	}

	idSet := make(map[int64]struct{}, len(rows))
	ids := make([]interface{}, 0, len(rows))
	for _, r := range rows {
		if _, ok := idSet[r.ID]; ok {
			continue
		}
		idSet[r.ID] = struct{}{}
		ids = append(ids, r.ID)
	}

	validated := make(map[int64]struct{}, len(ids))

	const chunk = 5000
	for i := 0; i < len(ids); i += chunk {
		j := i + chunk
		if j > len(ids) {
			j = len(ids)
		}
		part := ids[i:j]

		ph := make([]string, len(part))
		for k := range part {
			ph[k] = "?"
		}

		q := `
			SELECT DISTINCT address_id
			  FROM geo_address_validation
			 WHERE address_id IN (` + strings.Join(ph, ",") + `)
		`
		var rowsV []struct {
			AddressID int64 `orm:"column(address_id)"`
		}
		if _, err := o.Raw(q, part...).QueryRows(&rowsV); err != nil {
			return nil, nil, 0, err
		}
		for _, v := range rowsV {
			validated[v.AddressID] = struct{}{}
		}
	}

	filtered := make([]RowAddr, 0, len(rows))
	seenUUID := make(map[string]struct{}, len(rows))
	sendUUIDs := make([]string, 0, len(rows))

	for _, r := range rows {
		if _, ok := validated[r.ID]; ok {
			continue
		}
		u := strings.TrimSpace(r.UUID)
		if u == "" {
			continue
		}
		filtered = append(filtered, r)
		if _, ok := seenUUID[u]; !ok {
			seenUUID[u] = struct{}{}
			sendUUIDs = append(sendUUIDs, u)
		}
	}

	skipped := len(rows) - len(filtered)
	return filtered, sendUUIDs, skipped, nil
}
