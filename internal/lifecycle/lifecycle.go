package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/n24-x/pi/internal/piclock"
	"github.com/n24-x/pi/internal/pireflect"
	"github.com/n24-x/pi/pievent"
)

type Lifecycle struct {
	mu sync.Mutex

	state appState
	hooks []Hook

	numStarted   int
	startRecords HookRecords
	stopRecords  HookRecords
	runningHook  Hook

	clock  piclock.Time
	logger pievent.Logger
}

// New constructs a new [Lifecycle].
func New(logger pievent.Logger, clock piclock.Time) *Lifecycle {
	return &Lifecycle{
		logger: logger,
		clock:  clock,
	}
}

// AppendHook adds a Hook to the lifecycle.
func (l *Lifecycle) AppendHook(hook Hook) {
	// Save the caller's stack frame to report file/line number.
	if f := pireflect.CallerStack(2, 0); len(f) > 0 {
		hook.callerFrame = f[0]
	}

	l.hooks = append(l.hooks, hook)
}

// Start runs all OnStart hooks, returning immediately if it encounters an
// error.
func (l *Lifecycle) Start(ctx context.Context) error {
	if ctx == nil {
		return errors.New("called OnStart with nil context")
	}

	l.mu.Lock()
	if l.state != stopped {
		defer l.mu.Unlock()
		return fmt.Errorf("attempted to start lifecycle when in state: %v",
			l.state)
	}
	l.numStarted = 0
	l.state = starting

	l.startRecords = make(HookRecords, 0, len(l.hooks))
	l.mu.Unlock()

	returnState := incompleteStart
	defer func() {
		l.mu.Lock()
		l.state = returnState
		l.mu.Unlock()
	}()

	for _, hook := range l.hooks {
		// if ctx has cancelled, bail out of the loop.
		if err := ctx.Err(); err != nil {
			return err
		}

		if hook.OnStart != nil {
			l.mu.Lock()
			l.runningHook = hook
			l.mu.Unlock()

			elapsedTime, err := l.runStartHook(ctx, hook)
			if err != nil {
				return err
			}

			l.mu.Lock()
			l.startRecords = append(l.startRecords, HookRecord{
				CallerFrame: hook.callerFrame,
				Func:        hook.OnStart,
				ElapsedTime: elapsedTime,
			})
			l.mu.Unlock()
		}
		l.numStarted++
	}

	returnState = started
	return nil
}

func (l *Lifecycle) runStartHook(ctx context.Context, hook Hook) (
	elapsedTime time.Duration,
	err error,
) {
	funcName := hook.OnStartName
	if len(funcName) == 0 {
		funcName = pireflect.FuncName(hook.OnStart)
	}

	l.logger.LogEvent(&pievent.OnStartExecuting{
		CallerName:   hook.callerFrame.Function,
		FunctionName: funcName,
	})

	defer func() {
		l.logger.LogEvent(&pievent.OnStartExecuted{
			CallerName:   hook.callerFrame.Function,
			FunctionName: funcName,
			ElapsedTime:  elapsedTime,
			Err:          err,
		})
	}()

	begin := l.clock.Now()
	err = hook.OnStart(ctx)
	return l.clock.Since(begin), err
}

// Stop runs any OnStop hooks whose OnStart counterpart succeeded. OnStop
// hooks run in reverse order.
func (l *Lifecycle) Stop(ctx context.Context) error {
	if ctx == nil {
		return errors.New("called OnStop with nil context")
	}

	l.mu.Lock()
	if l.state != started && l.state != incompleteStart && l.state != starting {
		defer l.mu.Unlock()
		return nil
	}
	l.state = stopping
	l.mu.Unlock()

	defer func() {
		l.mu.Lock()
		l.state = stopped
		l.mu.Unlock()
	}()

	l.mu.Lock()
	l.stopRecords = make(HookRecords, 0, l.numStarted)
	// Take a snapshot of hook state to avoid races.
	allHooks := l.hooks[:]
	numStarted := l.numStarted
	l.mu.Unlock()

	// Run backward from last successful OnStart.
	var errs []error
	for ; numStarted > 0; numStarted-- {
		if err := ctx.Err(); err != nil {
			return err
		}
		hook := allHooks[numStarted-1]
		if hook.OnStop == nil {
			continue
		}

		l.mu.Lock()
		l.runningHook = hook
		l.mu.Unlock()

		elapsedTime, err := l.runStopHook(ctx, hook)
		if err != nil {
			// For best-effort cleanup, keep going after errors.
			errs = append(errs, err)
		}

		l.mu.Lock()
		l.stopRecords = append(l.stopRecords, HookRecord{
			CallerFrame: hook.callerFrame,
			Func:        hook.OnStop,
			ElapsedTime: elapsedTime,
		})
		l.mu.Unlock()
	}

	return errors.Join(errs...)
}

func (l *Lifecycle) runStopHook(ctx context.Context, hook Hook) (
	elapsedTime time.Duration,
	err error,
) {
	funcName := hook.OnStopName
	if len(funcName) == 0 {
		funcName = pireflect.FuncName(hook.OnStop)
	}

	l.logger.LogEvent(&pievent.OnStopExecuting{
		CallerName:   hook.callerFrame.Function,
		FunctionName: funcName,
	})
	defer func() {
		l.logger.LogEvent(&pievent.OnStopExecuted{
			CallerName:   hook.callerFrame.Function,
			FunctionName: funcName,
			ElapsedTime:  elapsedTime,
			Err:          err,
		})
	}()

	begin := l.clock.Now()
	err = hook.OnStop(ctx)
	return l.clock.Since(begin), err
}

// RunningHookCaller returns the name of the hook that was running when a Start/Stop
// hook timed out.
func (l *Lifecycle) RunningHookCaller() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.runningHook.callerFrame.Function
}
