package pireflect

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// Stack is a stack of call frames.
//
// Formatted with %v, the output is in a single-line, in the form,
//
//	foo/bar.Baz() (path/to/foo.go:42); bar/baz.Qux() (bar/baz/qux.go:12); ...
//
// Formatted with %+v, the output is in the form,
//
//	foo/bar.Baz()
//		path/to/foo.go:42
//	bar/baz.Qux()
//		bar/baz/qux.go:12
type Stack []Frame

// String returns a single-line, semi-colon representation of a Stack.
// For a list of strings where each represents one frame, use Strings.
// For a cleaner multi-line representation, use %+v.
func (fs Stack) String() string {
	return strings.Join(fs.Strings(), "; ")
}

// Strings returns a list of strings, each representing a frame in the stack.
// Each line will be in the form,
//
//	foo/bar.Baz() (path/to/foo.go:42)
func (fs Stack) Strings() []string {
	items := make([]string, len(fs))
	for i, f := range fs {
		items[i] = f.String()
	}
	return items
}

// Format implements fmt.Formatter to handle "%+v".
func (fs Stack) Format(w fmt.State, c rune) {
	if !w.Flag('+') {
		// Without %+v, fall back to String().
		io.WriteString(w, fs.String())
		return
	}

	for _, f := range fs {
		fmt.Fprintln(w, f.Function)
		fmt.Fprintf(w, "\t%v:%v\n", f.File, f.Line)
	}
}

// CallerName returns the name of the first caller in this stack that isn't
// owned by the Pi library.
func (fs Stack) CallerName() string {
	for _, f := range fs {
		if shouldIgnoreFrame(f) {
			continue
		}
		return f.Function
	}
	return "n/a"
}

const _defaultCallersDepth = 8

// CallerStack returns the call stack for the calling function, up to depth frames
// deep, skipping the provided number of frames, not including Callers itself.
//
// If zero, depth defaults to 8.
func CallerStack(skip int, depth int) Stack {
	if depth <= 0 {
		depth = _defaultCallersDepth
	}
	pcs := make([]uintptr, depth)

	// +2 to skip this frame and runtime.Callers.
	n := runtime.Callers(skip+2, pcs)
	pcs = pcs[:n] // truncate to number of frames actually read

	result := make([]Frame, 0, n)
	frames := runtime.CallersFrames(pcs)
	for f, more := frames.Next(); more; f, more = frames.Next() {
		result = append(result, Frame{
			Function: sanitizeFuncName(f.Function),
			File:     f.File,
			Line:     f.Line,
		})
	}

	return result
}
