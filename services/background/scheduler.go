package background

import (
	"context"
	"log"
	"runtime/debug"
	"time"
)

func RunDailyAt(ctx context.Context, hh, mm int, loc *time.Location, fn func(context.Context) error) {
	go func() {
		defer safeRecover("RunDailyAt")
		for {
			next := nextOccurrence(hh, mm, loc, time.Now())
			t := time.NewTimer(time.Until(next))
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				_ = fn(ctx)
			}
		}
	}()
}

func nextOccurrence(hh, mm int, loc *time.Location, ref time.Time) time.Time {
	now := ref.In(loc)
	next := time.Date(now.Year(), now.Month(), now.Day(), hh, mm, 0, 0, loc)
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

func safeRecover(where string) {
	if r := recover(); r != nil {
		log.Printf("[background][panic][%s]: %+v\n%s", where, r, debug.Stack())
	}
}
