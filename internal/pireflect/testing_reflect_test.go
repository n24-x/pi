package pireflect

import (
	"testing"

	"github.com/n24-x/goleak"
)

func TestCaller(t *testing.T) {
	t.Parallel()

	funcName := Caller()
	if funcName != "github.com/n24-x/pi/internal/pireflect.TestCaller" {
		t.Fatalf("Caller() = %q, want %q",
			funcName,
			"github.com/n24-x/pi/internal/pireflect.TestCaller",
		)
	}
}

func someFunc() {}

func TestFuncName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name:  "function",
			input: someFunc,
			want:  "github.com/n24-x/pi/internal/pireflect.someFunc()",
		},
		{
			name:  "not a function",
			input: 42,
			want:  "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			funcName := FuncName(tt.input)
			if tt.want != funcName {
				t.Fatalf("FuncName(%v) = %q, want %q", tt.input, funcName, tt.want)
			}
		})
	}
}

func TestSanitizeFuncNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"url encoding",
			"github.com/n24-x/pi/sample%2egit/someFunc",
			"github.com/n24-x/pi/sample.git/someFunc",
		},
		{
			"vendor removal",
			"github.com/n24-x/pi/vendor/github.com/some/lib.SomeFunc",
			"vendor/github.com/some/lib.SomeFunc",
		},
		{
			"package happens to be named vendor is untouched",
			"github.com/n24-x/pi/foovendor/someFunc",
			"github.com/n24-x/pi/foovendor/someFunc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			funcName := sanitizeFuncName(tt.input)
			if tt.want != funcName {
				t.Fatalf("sanitizeFuncName(%q) = %q, want %q",
					tt.input, funcName, tt.want)
			}

		})
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
