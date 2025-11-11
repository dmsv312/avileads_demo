package background

import (
	"avileads-web/services/booking"
	"avileads-web/services/jobs"
	"context"
	"github.com/astaxie/beego/orm"
	"time"
)

type NewOrm func() orm.Ormer

func Start(ctx context.Context, newOrm NewOrm, d *booking.DadataService) {
	go func() {
		defer safeRecover("background.Start")
		msk, _ := time.LoadLocation("Europe/Moscow")

		jobs.StartWorker(ctx, newOrm(), d)

		RunDailyAt(ctx, 23, 00, msk, func(c context.Context) error {
			return jobs.PlanDailySimple(c, jobs.NewOrm(newOrm), d)
		})

		RunDailyAt(ctx, 3, 00, msk, func(c context.Context) error {
			return jobs.CleanupBookings(c, newOrm())
		})
	}()
}
