//go:build go1.27 && !go1.28

package pireflect

import (
	"strings"
	"testing"
)

// Files are split by Go release because the compiler-generated name of the
// innermost nested function literal changes between versions
// When a newer Go changes it again, add testing_stack_go1.<n>_test.go and
// bound this file's tag: //go:build go1.27 && !go1.28

func TestDeepStack(t *testing.T) {
	t.Run("nest", func(t *testing.T) {
		// Introduce a few frames.
		frames := func() []Frame {
			return func() []Frame {
				return CallerStack(0, 0)
			}()
		}()

		if len(frames) <= 3 {
			t.Fatalf("CallerStack(0, 0) returned %d frames, want at least 4", len(frames))
		}
		for i, name := range []string{"func1.1.1", "func1.1", "func1"} {
			f := frames[i]
			want := "github.com/n24-x/pi/internal/pireflect.TestDeepStack." + name
			if f.Function != want {
				t.Fatalf("frame function = %q, want %q", f.Function, want)
			}
			if !strings.Contains(f.File, "internal/pireflect/testing_stack_go1.27_test.go") {
				t.Fatalf("frame file = %q, want it to contain %q",
					f.File, "internal/pireflect/testing_stack_go1.27_test.go")
			}
			if f.Line == 0 {
				t.Fatal("frame line = 0, want non-zero")
			}
		}
	})
}
