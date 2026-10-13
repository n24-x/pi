package pireflect

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestStack(t *testing.T) {
	t.Parallel()

	// NOTE:
	// We don't assert the length of the stack because we cannot make
	// guarantees about how many frames the test runner is allowed to
	// introduce.

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		frames := CallerStack(0, 0)
		if len(frames) == 0 {
			t.Fatal("CallerStack(0, 0) returned no frames")
		}
		f := frames[0]
		if want := "github.com/n24-x/pi/internal/pireflect.TestStack.func1"; f.Function != want {
			t.Fatalf("frame function = %q, want %q", f.Function, want)
		}
		if !strings.Contains(f.File, "internal/pireflect/testing_stack_test.go") {
			t.Fatalf("frame file = %q, want it to contain %q",
				f.File, "internal/pireflect/testing_stack_test.go")
		}
		if f.Line == 0 {
			t.Fatal("frame line = 0, want non-zero")
		}
	})

	t.Run("skip", func(t *testing.T) {
		t.Parallel()

		// Introduce a few frames and skip 2.
		frames := func() []Frame {
			return func() []Frame {
				return CallerStack(2, 0)
			}()
		}()

		if len(frames) == 0 {
			t.Fatal("CallerStack(2, 0) returned no frames")
		}

		f := frames[0]
		want := "github.com/n24-x/pi/internal/pireflect.TestStack.func2"
		if f.Function != want {
			t.Fatalf("frame function = %q, want %q", f.Function, want)
		}

		if !strings.Contains(f.File, "internal/pireflect/testing_stack_test.go") {
			t.Fatalf("frame file = %q, want it to contain %q",
				f.File, "internal/pireflect/testing_stack_test.go")
		}

		if f.Line == 0 {
			t.Fatal("frame line = 0, want non-zero")
		}
	})
}

func TestStackCallerName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input Stack
		want  string
	}{
		{name: "empty", want: "n/a"},
		{
			name: "skip Pi components",
			input: Stack{
				{
					Function: "github.com/n24-x/pi.Foo()",
					File:     "github.com/n24-x/pi/foo.go",
				},
				{
					Function: "foo/bar.Baz()",
					File:     "foo/bar/baz.go",
				},
			},
			want: "foo/bar.Baz()",
		},
		{
			name: "skip Pi in wrong directory",
			input: Stack{
				{
					Function: "github.com/n24-x/pi.Foo()",
					File:     "pi/foo.go",
				},
				{
					Function: "foo/bar.Baz()",
					File:     "foo/bar/baz.go",
				},
			},
			want: "foo/bar.Baz()",
		},
		{
			name: "skip Pi subpackage",
			input: Stack{
				{
					Function: "github.com/n24-x/pi/internal/pireflect.Foo()",
					File:     "pi/internal/pireflect/foo.go",
				},
				{
					Function: "foo/bar.Baz()",
					File:     "foo/bar/baz.go",
				},
			},
			want: "foo/bar.Baz()",
		},
		{
			name: "don't skip Pi tests",
			input: Stack{
				{
					Function: "some/thing.Foo()",
					File:     "github.com/n24-x/pi/foo_test.go",
				},
			},
			want: "some/thing.Foo()",
		},
		{
			name: "don't skip Pi prefix",
			input: Stack{
				{
					Function: "github.com/n24-x/pifoo.Bar()",
					File:     "github.com/n24-x/pifoo/bar.go",
				},
			},
			want: "github.com/n24-x/pifoo.Bar()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.input.CallerName()
			if got != tt.want {
				t.Fatalf("CallerName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFrameString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input Frame
		want  string
	}{
		{
			name:  "zero",
			input: Frame{},
			want:  "unknown",
		},
		{
			name:  "file and line",
			input: Frame{File: "foo.go", Line: 42},
			want:  "(foo.go:42)",
		},
		{
			name:  "file only",
			input: Frame{File: "foo.go"},
			want:  "(foo.go)",
		},
		{
			name:  "function only",
			input: Frame{Function: "foo"},
			want:  "foo",
		},
		{
			name:  "function and file",
			input: Frame{Function: "foo", File: "bar.go"},
			want:  "foo (bar.go)",
		},
		{
			name:  "function and line",
			input: Frame{Function: "foo", Line: 42},
			want:  "foo", // line without file is meaningless
		},
		{
			name:  "function, file, and line",
			input: Frame{Function: "foo", File: "bar.go", Line: 42},
			want:  "foo (bar.go:42)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.input.String(); got != tt.want {
				t.Fatalf("Frame.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStackFormat(t *testing.T) {
	t.Parallel()

	stack := Stack{
		{
			Function: "path/to/module.SomeFunction()",
			File:     "path/to/file.go",
			Line:     42,
		},
		{
			Function: "path/to/another/module.AnotherFunction()",
			File:     "path/to/another/file.go",
			Line:     12,
		},
	}

	t.Run("single line", func(t *testing.T) {
		t.Parallel()

		want := "path/to/module.SomeFunction() (path/to/file.go:42); " +
			"path/to/another/module.AnotherFunction() (path/to/another/file.go:12)"
		if got := fmt.Sprintf("%v", stack); got != want {
			t.Fatalf("Stack %%v = %q, want %q", got, want)
		}
	})

	t.Run("multi line", func(t *testing.T) {
		t.Parallel()

		want := `path/to/module.SomeFunction()
	path/to/file.go:42
path/to/another/module.AnotherFunction()
	path/to/another/file.go:12
`
		if got := fmt.Sprintf("%+v", stack); got != want {
			t.Fatalf("Stack %%+v = %q, want %q", got, want)
		}
	})

	t.Run("strings", func(t *testing.T) {
		t.Parallel()

		want := []string{
			"path/to/module.SomeFunction() (path/to/file.go:42)",
			"path/to/another/module.AnotherFunction() (path/to/another/file.go:12)",
		}
		if got := stack.Strings(); !slices.Equal(got, want) {
			t.Fatalf("Stack.Strings() = %q, want %q", got, want)
		}
	})
}
