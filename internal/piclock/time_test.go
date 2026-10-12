package piclock

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestSystemTime(t *testing.T) {
	clock := System
	testTime(t, System, clock.Sleep)
}

func TestMockTime(t *testing.T) {
	clock := NewMock()
	testTime(t, clock, clock.Advance)
}

func TestMock_AdvanceNegative(t *testing.T) {
	clock := NewMock()

	defer func() {
		if r := recover(); r == nil {
			t.Error("advance with negative duration must panic")
		}
	}()
	clock.Advance(-1)
}

func TestMock_Sleep(t *testing.T) {
	clock := NewMock()

	ch := make(chan struct{})
	go func() {
		clock.Sleep(2 * time.Millisecond)
		close(ch)
	}()

	// We cannot advance time until we're certain
	// that the Sleep call has started waiting.
	// Otherwise, we'll advance that one millisecond,
	// and then the Sleep will start waiting for another Advance,
	// which will never come.
	//
	// AwaitScheduled will block until there is at least one
	// scheduled event.
	clock.AwaitScheduled(1)

	// Advance only one millisecond, the Sleep should not return.
	clock.Advance(1 * time.Millisecond)
	select {
	case <-ch:
		t.Fatal("sleep should not have returned")
	case <-time.After(1 * time.Millisecond):
		// expected: Sleep still blocked
	}

	clock.Advance(1 * time.Millisecond)
	select {
	case <-ch:
		// expected: Sleep returned
	case <-time.After(10 * time.Millisecond):
		t.Fatal("expected Sleep to return")
	}
}

func TestMock_ManySleepers(t *testing.T) {
	const N = 200

	clock := NewMock()

	var wg sync.WaitGroup
	wg.Add(N)
	for range N {
		go func() {
			defer wg.Done()

			clock.Sleep(1 * time.Millisecond)
		}()
	}

	clock.AwaitScheduled(N)
	clock.Advance(1 * time.Millisecond)

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		// expected: all sleepers are done
	case <-time.After(10 * time.Millisecond):
		t.Fatal("expected all sleepers to be done")
	}
}

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

		clock.Sleep(1 * time.Millisecond) // add a waiter; see time_mock.go L131
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
			t.Fatal("timed out waiting for context to be done")
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

	t.Run("WithTimeout/Cancel", func(t *testing.T) {
		ctx, cancel := clock.WithTimeout(context.Background(), 1*time.Millisecond)
		cancel()

		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err == nil {
				t.Error("done context must error")
			}
			if !errors.Is(err, context.Canceled) {
				t.Error("context must have been canceled")
			}

		case <-time.After(10 * time.Millisecond):
			t.Fatal("timed out waiting for context to be done")
		}
	})
}
