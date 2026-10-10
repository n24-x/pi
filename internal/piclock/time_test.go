package piclock

import (
	"context"
	"errors"
	"testing"
	"time"
)

// func TestSystemTime(t *testing.T) {
// 	time := System
// }

func testTime(t *testing.T, clock Time, advanceClock func(d time.Duration)) {
	now := clock.Now()
	if now.IsZero() {
		t.Error("time must not return zero time")
	}

	t.Run("Since", func(t *testing.T) {
		advanceClock(1 * time.Millisecond)
		if clock.Since(now) == 0 {
			t.Error("time must have advanced")
		}
	})

	t.Run("Sleep", func(t *testing.T) {
		start := clock.Now()
		go func() {
			// For the mock clock, there's a chance that advanceClock will be
			// too fast and the Sleep will block forever, waiting for
			// another advanceClock. The mock clock provides
			// AwaitScheduled to help with this.
			//
			// So that function is not available on the system clock,
			// we'll use upcasting to check for it.
			if awaiter, ok := clock.(interface{ AwaitScheduled(int) }); ok {
				awaiter.AwaitScheduled(1)
			}

			advanceClock(1 * time.Millisecond)
		}()

		clock.Sleep(1 * time.Millisecond)
		if clock.Since(start) == 0 {
			t.Error("time must have advanced")
		}
	})

	t.Run("WithTimeout", func(t *testing.T) {
		ctx, cancel := clock.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		t.Run("Deadline", func(t *testing.T) {
			dl, ok := ctx.Deadline()
			if !ok {
				t.Error("must have a deadline")
			}
			if !dl.After(now) {
				t.Error("deadline must be in the future")
			}
		})

		advanceClock(1 * time.Millisecond)

		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err == nil {
				t.Error("done context must error")
			}
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Error("context must have exceeded its deadline")
			}
		case <-time.After(10 * time.Millisecond):
			t.Fatalf("expected context to be done")
		}
	})

	t.Run("WithTimeout/Value", func(t *testing.T) {
		type contextKey string
		key := contextKey("foo")

		ctx1 := context.WithValue(context.Background(), key, "bar")

		ctx2, cancel := clock.WithTimeout(ctx1, 1*time.Millisecond)
		defer cancel()

		if ctx2.Value(key) != "bar" {
			t.Error("value must be preserved")
		}
	})
}
