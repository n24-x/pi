package pitime

import (
	"context"
	"sync"
	"time"
)

type deadlineCtx struct {
	inner       context.Context
	cancelInner func()

	done     chan struct{}
	deadline time.Time

	mu  sync.Mutex // guards err; the rest is immutable
	err error
}

var _ context.Context = (*deadlineCtx)(nil)

func (c *deadlineCtx) Deadline() (deadline time.Time, ok bool) { return c.deadline, true }
func (c *deadlineCtx) Done() <-chan struct{}                   { return c.done }
func (c *deadlineCtx) Value(key any) any                       { return c.inner.Value(key) }

func (c *deadlineCtx) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func (c *deadlineCtx) cancel(err error) {
	c.mu.Lock()
	if c.err == nil {
		c.err = err
		close(c.done)
		c.cancelInner()
	}
	c.mu.Unlock()
}
