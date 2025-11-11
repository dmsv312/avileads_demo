package jobs

import (
	"avileads-web/services/booking"
	"context"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"log"
	"runtime/debug"
	"time"
)

func StartWorker(ctx context.Context, o orm.Ormer, dadata *booking.DadataService) {
	go func() {
		defer safeRecover("jobs.StartWorker")
		loop(ctx, o, dadata)
	}()
}

func loop(ctx context.Context, o orm.Ormer, dadata *booking.DadataService) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		var jobID int64
		var rawPayload string

		err := o.Raw(`
			UPDATE geo_background_job
			   SET status='running', started_at=now()
			 WHERE id = (
			       SELECT id FROM geo_background_job
			        WHERE status='pending'
			        ORDER BY id
			        FOR UPDATE SKIP LOCKED
			        LIMIT 1
			   )
			 RETURNING id, payload
		`).QueryRow(&jobID, &rawPayload)
		if err != nil || jobID == 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
			continue
		}
		payload := []byte(rawPayload)

		if err := handleBooking(o, jobID, payload, dadata); err != nil {
			o.Raw(`UPDATE geo_background_job
			          SET status='error', error=?, finished_at=now()
			        WHERE id=?`, err.Error(), jobID).Exec()
		} else {
			o.Raw(`UPDATE geo_background_job
			          SET status='done',  finished_at=now()
			        WHERE id=?`, jobID).Exec()
		}
	}
}

type bookingJobPayload struct {
	UserID   int                   `json:"user_id"`
	Bookings []booking.BookingItem `json:"bookings"`
}

func handleBooking(o orm.Ormer, jobID int64, raw []byte, d *booking.DadataService) error {
	var p bookingJobPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return err
	}
	return booking.ProcessAndSave(o, p.UserID, p.Bookings, d)
}

func safeRecover(where string) {
	if r := recover(); r != nil {
		log.Printf("[panic recovered] %s: %v\n%s", where, r, debug.Stack())
	}
}
