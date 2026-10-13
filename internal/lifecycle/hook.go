package lifecycle

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/n24-x/pi/internal/pireflect"
)

type Hook struct {
	OnStart     func(context.Context) error
	OnStartName string

	OnStop     func(context.Context) error
	OnStopName string

	callerFrame pireflect.Frame
}

// HookRecord keeps track of each Hook's execution time, the caller that appended the Hook, and function that ran as the Hook.
type HookRecord struct {
	CallerFrame pireflect.Frame             // stack frame of the caller
	Func        func(context.Context) error // function that ran as sanitized name
	ElapsedTime time.Duration               // how long the hook ran
}

// HookRecords is a Stringer wrapper of HookRecord slice.
type HookRecords []HookRecord

func (rs HookRecords) Len() int {
	return len(rs)
}

func (rs HookRecords) Less(i, j int) bool {
	// Sort by time, greater ones at top.
	return rs[i].ElapsedTime > rs[j].ElapsedTime
}

func (rs HookRecords) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// Used for logging startup errors.
func (rs HookRecords) String() string {
	var b strings.Builder
	for _, r := range rs {
		fmt.Fprintf(&b, "%s took %v from %s",
			pireflect.FuncName(r.Func), r.ElapsedTime, r.CallerFrame)
	}
	return b.String()
}

// Format implements fmt.Formatter to handle "%+v".
func (rs HookRecords) Format(w fmt.State, c rune) {
	if !w.Flag('+') {
		// Without %+v, fallback to String().
		io.WriteString(w, rs.String())
		return
	}

	for _, r := range rs {
		fmt.Fprintf(w, "\n%s took %v from:\n\t%+v",
			pireflect.FuncName(r.Func),
			r.ElapsedTime,
			r.CallerFrame)
	}

	fmt.Fprintf(w, "\n")
}
