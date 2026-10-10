package piclock

import (
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

}
