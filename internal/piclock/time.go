package piclock

import (
	"context"
	"time"
)

// Time defines how Pi accesses time.
// We keep the interface pretty minimal.
type Time interface {
	Now() time.Time
	Since(time.Time) time.Duration
	Sleep(time.Duration)
	WithTimeout(context.Context, time.Duration) (context.Context, context.CancelFunc)
}

// System is the default implementation of [Time] based on real time.
var System Time = systemTime{}

type systemTime struct{}

func (st systemTime) Now() time.Time {
	return time.Now()
}

func (st systemTime) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (st systemTime) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (st systemTime) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}
