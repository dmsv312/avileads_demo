package jobs

import (
	"context"
	"github.com/astaxie/beego/orm"
	"log"
	"time"
)

const cleanupDeactivateBookingsSQL = `
UPDATE public.geo_address_booking b
SET    is_active = FALSE
WHERE  b.is_active = TRUE
  AND  b.end_date IS NOT NULL
  AND  b.end_date <= NOW();
`

const cleanupAddressStatusSQL = `
UPDATE public.geo_address_status s
SET    status_id = 1
FROM (
  SELECT DISTINCT ON (address_id)
         address_id, end_date
  FROM   public.geo_address_booking
  WHERE  end_date IS NOT NULL
  ORDER  BY address_id, end_date DESC
) last
WHERE  last.address_id = s.address_id
  AND   s.status_id <> 1
  AND   last.end_date <= NOW();
`

func CleanupBookings(ctx context.Context, o orm.Ormer) error {
	start := time.Now()
	log.Println("[cleanup] start")

	res1, err := o.Raw(cleanupDeactivateBookingsSQL).Exec()
	if err != nil {
		log.Printf("[cleanup] SQL error (deactivate bookings): %v", err)
		return err
	}
	affectedBookings, _ := res1.RowsAffected()
	log.Printf("[cleanup] bookings deactivated: %d", affectedBookings)

	res2, err := o.Raw(cleanupAddressStatusSQL).Exec()
	if err != nil {
		log.Printf("[cleanup] SQL error (update address status): %v", err)
		return err
	}
	affectedStatuses, _ := res2.RowsAffected()

	log.Printf("[cleanup] done: affected_status=%d, took=%s", affectedStatuses, time.Since(start))
	return nil
}
