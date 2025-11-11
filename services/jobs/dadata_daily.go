package jobs

import (
	"avileads-web/services/booking"
	"avileads-web/services/validation"
	"context"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"log"
)

const (
	defaultMinBookings = 100
	defaultRPS         = 17
	defaultWorkers     = 12
	defaultLogEvery    = 1000
)

type planRow struct {
	CityID        int64   `orm:"column(city_id)"`
	CityName      string  `orm:"column(city_name)"`
	RegionName    string  `orm:"column(region_name)"`
	BookingsCount int64   `orm:"column(bookings_count)"`
	PendingCount  int64   `orm:"column(pending_count)"`
	SharePct      float64 `orm:"column(share_pct)"`
	BaseAlloc     int     `orm:"column(base_alloc)"`
	ExtraAlloc    int     `orm:"column(extra_alloc)"`
	Allocated     int     `orm:"column(allocated)"`
}

const planSQL = `
WITH
params AS (
   SELECT ?::int AS budget, ?::int AS min_bookings
),
bookings AS (
   SELECT
       c.id   AS city_id,
       c.name AS city_name,
       r.name AS region_name,
       COUNT(*)::bigint AS bookings_count
   FROM public.geo_address_booking b
   JOIN avileads_geo.geo_address a ON a.id = b.address_id
   JOIN avileads_geo.geo_city    c ON c.id = a.city_id
   JOIN avileads_geo.geo_region  r ON r.id = c.region_id
   WHERE c.id <> 1405113 AND c.id <> 1414662
   GROUP BY c.id, c.name, r.name
),
visible AS (
   SELECT b.* FROM bookings b, params p
   WHERE b.bookings_count >= p.min_bookings
),
pending AS (
   SELECT a.city_id, COUNT(*)::bigint AS pending_count
   FROM avileads_geo.geo_address a
   JOIN avileads_geo.geo_city   c ON c.id = a.city_id
   JOIN avileads_geo.geo_region r ON r.id = c.region_id
   WHERE a.location IS NULL
     AND a.uuid IS NOT NULL
     AND c.id <> 1405113 AND c.id <> 1414662
     AND NOT EXISTS (
           SELECT 1 FROM public.geo_address_validation v
           WHERE v.address_id = a.id
             AND v.is_match   = TRUE
             AND v.qc_geo     <= 1
             AND v.qc_house   = 0
     )
   GROUP BY a.city_id
),
eligible AS (
   SELECT v.city_id, v.city_name, v.region_name, v.bookings_count,
          COALESCE(p.pending_count, 0) AS pending_count
   FROM visible v JOIN pending p ON p.city_id = v.city_id
   WHERE COALESCE(p.pending_count, 0) > 0
),
weights AS (
   SELECT e.*, SUM(e.bookings_count) OVER ()::numeric AS total_weight
   FROM eligible e
),
base AS (
   SELECT w.*,
          CASE
            WHEN NULLIF(w.total_weight, 0) IS NULL THEN 0
            ELSE LEAST(
                   FLOOR( ((SELECT budget FROM params)::numeric * w.bookings_count::numeric)
                          / NULLIF(w.total_weight, 0) )::int,
                   w.pending_count::int
                 )
          END AS base_alloc
   FROM weights w
),
sums AS ( SELECT SUM(base_alloc)::int AS base_sum,
                 SUM(pending_count)::int AS total_pending
          FROM base),
leftover AS (
   SELECT GREATEST( LEAST((SELECT budget FROM params), s.total_pending) - s.base_sum, 0 )::int AS left1
   FROM sums s
),
room AS (
   SELECT b.*, (b.pending_count::int - b.base_alloc::int) AS room_left
   FROM base b
),
room_weights AS (
   SELECT r.*, SUM(CASE WHEN r.room_left > 0 THEN r.bookings_count ELSE 0 END) OVER ()::numeric AS total_w2,
          (SELECT left1 FROM leftover) AS left1
   FROM room r
),
extra_main AS (
   SELECT rw.*,
          CASE
            WHEN rw.left1 = 0 OR NULLIF(rw.total_w2, 0) IS NULL OR rw.room_left <= 0 THEN 0
            ELSE LEAST(
                   FLOOR( (rw.left1::numeric * rw.bookings_count::numeric) / NULLIF(rw.total_w2, 0) )::int,
                   rw.room_left
                 )
          END AS extra_floor
   FROM room_weights rw
),
extra_sums AS ( SELECT SUM(extra_floor)::int AS extra_sum FROM extra_main ),
leftover2 AS (
   SELECT GREATEST( (SELECT left1 FROM leftover) - (SELECT extra_sum FROM extra_sums), 0 )::int AS left2
),
ranked AS (
   SELECT em.*,
          CASE
            WHEN em.left1 = 0 OR NULLIF(em.total_w2, 0) IS NULL OR (em.room_left - em.extra_floor) <= 0
                 THEN 0::numeric
            ELSE ((em.left1::numeric * em.bookings_count::numeric) / NULLIF(em.total_w2, 0))
                 - FLOOR((em.left1::numeric * em.bookings_count::numeric) / NULLIF(em.total_w2, 0))
          END AS frac_part,
          ROW_NUMBER() OVER (
              ORDER BY
                CASE WHEN (em.room_left - em.extra_floor) > 0 THEN
                     (((em.left1::numeric * em.bookings_count::numeric) / NULLIF(em.total_w2, 0))
                      - FLOOR((em.left1::numeric * em.bookings_count::numeric) / NULLIF(em.total_w2, 0)))
                ELSE -1::numeric END DESC,
                em.bookings_count DESC,
                em.city_id
          ) AS rn
   FROM extra_main em
)
SELECT
   r.city_id, r.city_name, r.region_name,
   r.bookings_count, r.pending_count,
   ROUND(
       COALESCE(
           r.bookings_count::numeric * 100
           / NULLIF( (SELECT SUM(bookings_count)::numeric FROM eligible), 0 ),
           0
       ), 2
   ) AS share_pct,
   r.base_alloc AS base_alloc,
   (r.extra_floor
    + CASE WHEN r.rn <= (SELECT left2 FROM leftover2)
             AND (r.room_left - r.extra_floor) > 0
           THEN 1 ELSE 0 END) AS extra_alloc,
   (r.base_alloc
    + r.extra_floor
    + CASE WHEN r.rn <= (SELECT left2 FROM leftover2)
             AND (r.room_left - r.extra_floor) > 0
           THEN 1 ELSE 0 END) AS allocated
FROM ranked r
WHERE (r.base_alloc + r.extra_floor
      + CASE WHEN r.rn <= (SELECT left2 FROM leftover2)
               AND (r.room_left - r.extra_floor) > 0
             THEN 1 ELSE 0 END) > 0
ORDER BY allocated DESC, r.bookings_count DESC, r.city_name;`

func getPlan(o orm.Ormer, budget, minBookings int) ([]planRow, error) {
	var rows []planRow
	_, err := o.Raw(planSQL, budget, minBookings).QueryRows(&rows)
	return rows, err
}

func fetchRows(o orm.Ormer, cityID int64, lim int) ([]validation.RowAddr, error) {
	const q = `
SELECT ga.id, ga.uuid, ga.city_id
FROM   avileads_geo.geo_address ga
JOIN   avileads_geo.geo_city   gc ON gc.id = ga.city_id
WHERE  ga.location IS NULL
AND    ga.uuid IS NOT NULL
AND    ($1 = 0 OR gc.id = $1)
AND NOT EXISTS (
      SELECT 1 FROM public.geo_address_validation v
      WHERE  v.address_id = ga.id
         AND v.is_match   = TRUE
)
ORDER BY ga.id
LIMIT  $2`
	var rows []validation.RowAddr
	_, err := o.Raw(q, cityID, lim).QueryRows(&rows)
	return rows, err
}

type NewOrm func() orm.Ormer

func PlanDaily(
	ctx context.Context,
	newOrm NewOrm,
	d validation.DadataClient,
	minBookings, rps, workers int,
	logEvery int,
) error {
	o := newOrm()

	remaining, err := dailyRemaining(o)
	if err != nil {
		return err
	}

	plan, err := getPlan(o, remaining, minBookings)
	if err != nil {
		return err
	}
	if len(plan) == 0 {
		log.Print("[daily] план пуст (нет pending) — выходим")
		return nil
	}

	rowsByCity := make(map[int64][]validation.RowAddr, len(plan))
	allIDs := make([]string, 0, remaining)
	totalPlanned := 0

	for _, p := range plan {
		if p.Allocated <= 0 {
			continue
		}
		rows, ferr := fetchRows(o, p.CityID, p.Allocated)
		if ferr != nil {
			log.Printf("[daily][city=%d] fetchRows: %v", p.CityID, ferr)
			continue
		}
		if len(rows) == 0 {
			continue
		}
		rowsByCity[p.CityID] = rows
		for _, r := range rows {
			allIDs = append(allIDs, r.UUID)
		}
		totalPlanned += len(rows)

		log.Printf("[daily][plan] %s/%s: allocated=%d pending=%d got=%d",
			p.RegionName, p.CityName, p.Allocated, p.PendingCount, len(rows))
	}

	if len(allIDs) == 0 {
		log.Print("[daily] нечего валидировать — выходим")
		return nil
	}

	log.Printf("[daily] pipeline start: total=%d, rps=%d, workers=%d", len(allIDs), rps, workers)
	pout := validation.RunPipeline(ctx, d, allIDs, validation.DefaultPipelineOptions)
	if len(pout.Errors) > 0 {
		log.Printf("[daily] pipeline finished with errors: %d (sent=%d)", len(pout.Errors), pout.TotalSent)
	}

	applied := 0
	for cityID, rows := range rowsByCity {
		oc := newOrm()
		if err := booking.ApplyUpdates(oc, rows, pout.Responses, logEvery); err != nil {
			log.Printf("[daily][city=%d] apply updates: %v", cityID, err)
			continue
		}
		applied += len(rows)
	}

	log.Printf("ALL DONE (multi-city job): planned=%d, sent=%d, applied_rows=%d", totalPlanned, pout.TotalSent, applied)
	return nil
}

func PlanDailySimple(ctx context.Context, newOrm NewOrm, d *booking.DadataService) error {
	return PlanDaily(ctx, newOrm, d, defaultMinBookings, defaultRPS, defaultWorkers, defaultLogEvery)
}

func dailyRemaining(o orm.Ormer) (int, error) {
	const q = `
	WITH bounds AS (
	  SELECT
		((now() AT TIME ZONE 'Europe/Moscow')::date::timestamp AT TIME ZONE 'Europe/Moscow') AS msk_start_utc,
		(((now() AT TIME ZONE 'Europe/Moscow')::date + 1)::timestamp AT TIME ZONE 'Europe/Moscow') AS msk_next_utc
	)
	SELECT GREATEST(49990 - COUNT(*), 0) AS remaining
	FROM public.geo_address_validation v, bounds b
	WHERE v.created_at >= b.msk_start_utc
	  AND v.created_at <  b.msk_next_utc;`
	var rem int
	if err := o.Raw(q).QueryRow(&rem); err != nil {
		return 0, err
	}
	return rem, nil
}
