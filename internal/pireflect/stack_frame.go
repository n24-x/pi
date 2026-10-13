package pireflect

import (
	"fmt"
	"strings"
)

// Frame holds information about a single frame in the call stack.
type Frame struct {
	// Unique, package path-qualified name for the function of this call
	// frame.
	Function string

	// File and line number of our location in the frame.
	//
	// Note that the line number does not refer to where the function was
	// defined but where in the function the next call was made.
	File string
	Line int
}

func (f Frame) String() string {
	// This takes the following forms.
	//  (path/to/file.go)
	//  (path/to/file.go:42)
	//  path/to/package.MyFunction
	//  path/to/package.MyFunction (path/to/file.go)
	//  path/to/package.MyFunction (path/to/file.go:42)

	var sb strings.Builder
	sb.WriteString(f.Function)
	if len(f.File) > 0 {
		if sb.Len() > 0 {
			sb.WriteRune(' ')
		}
		fmt.Fprintf(&sb, "(%v", f.File)
		if f.Line > 0 {
			fmt.Fprintf(&sb, ":%d", f.Line)
		}
		sb.WriteRune(')')
	}

	if sb.Len() == 0 {
		return "unknown"
	}

	return sb.String()
}

// Ascend the call stack until we leave the Pi production code. This allows us
// to avoid hard-coding a frame skip, which makes this code work well even
// when it's wrapped.
func shouldIgnoreFrame(f Frame) bool {
	// Treat test files as leafs.
	if strings.Contains(f.File, "_test.go") {
		return false
	}

	// The unique, fully-qualified name for all functions begins with
	// "{{importPath}}.". We'll ignore Pi and its subpackages.
	s := strings.TrimPrefix(f.Function, "github.com/n24-x/pi")
	if len(s) > 0 && s[0] == '.' || s[0] == '/' {
		// We want to match,
		//   github.com/n24-x/pi.Foo
		//   github.com/n24-x/pi/something.Foo
		// But not, github.com/n24-x/pifoo
		return true
	}

	return false
}
