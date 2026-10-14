package pievent

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// richError prints a different output when formatted with %+v vs %v.
type richError struct{}

func (e *richError) Error() string { return "plain error" }

func (e *richError) Format(w fmt.State, c rune) {
	if w.Flag('+') && c == 'v' {
		// Format differently for %+v.
		io.WriteString(w, "rich error")
	} else {
		io.WriteString(w, e.Error())
	}
}

func TestConsoleLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input Event
		want  string
	}{
		{
			name: "OnStart executing",
			input: &OnStartExecuting{
				FunctionName: "hook.onStart",
				CallerName:   "bytes.NewBuffer",
			},
			want: "[Pi] HOOK OnStart		hook.onStart executing (caller: bytes.NewBuffer)\n",
		},
		{
			name: "OnStopExecuting",
			input: &OnStopExecuting{
				FunctionName: "hook.onStop1",
				CallerName:   "bytes.NewBuffer",
			},
			want: "[Pi] HOOK OnStop		hook.onStop1 executing (caller: bytes.NewBuffer)\n",
		},
		{
			name: "OnStopExecutedError",
			input: &OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          fmt.Errorf("some error"),
			},
			want: "[Pi] HOOK OnStop		hook.onStart1 called by bytes.NewBuffer failed in 0s: some error\n",
		},
		{
			name: "OnStopExecutedError/rich error",
			input: &OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          &richError{},
			},
			want: "[Pi] HOOK OnStop		hook.onStart1 called by bytes.NewBuffer failed in 0s: rich error\n",
		},
		{
			name: "OnStopExecuted",
			input: &OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				ElapsedTime:  time.Millisecond * 3,
			},
			want: "[Pi] HOOK OnStop		hook.onStart1 called by bytes.NewBuffer ran successfully in 3ms\n",
		},
		{
			name: "OnStartExecutedError",
			input: &OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          fmt.Errorf("some error"),
			},
			want: "[Pi] HOOK OnStart		hook.onStart1 called by bytes.NewBuffer failed in 0s: some error\n",
		},
		{
			name: "OnStartExecutedError/rich error",
			input: &OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          &richError{},
			},
			want: "[Pi] HOOK OnStart		hook.onStart1 called by bytes.NewBuffer failed in 0s: rich error\n",
		},
		{
			name: "OnStartExecuted",
			input: &OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				ElapsedTime:  time.Millisecond * 3,
			},
			want: "[Pi] HOOK OnStart		hook.onStart1 called by bytes.NewBuffer ran successfully in 3ms\n",
		},
		{
			name:  "ResourceAddError",
			input: &ResourceAdded{Err: errors.New("some error")},
			want:  "[Pi] Error after options were applied: some error\n",
		},
		{
			name:  "ResourceAddError/rich error",
			input: &ResourceAdded{Err: &richError{}},
			want:  "[Pi] Error after options were applied: rich error\n",
		},
		{
			name: "Supplied",
			input: &Supplied{
				TypeName:    "*bytes.Buffer",
				StackTrace:  []string{"main.main", "runtime.main"},
				ModuleTrace: []string{"main.main"},
			},
			want: "[Pi] SUPPLY	*bytes.Buffer\n",
		},
		{
			name: "Supplied with module",
			input: &Supplied{
				TypeName:    "*bytes.Buffer",
				ModuleName:  "myModule",
				StackTrace:  []string{"main.main", "runtime.main"},
				ModuleTrace: []string{"main.main", "mypackage.GetMyModule"},
			},
			want: "[Pi] SUPPLY	*bytes.Buffer from module \"myModule\"\n",
		},
		{
			name:  "SuppliedError",
			input: &Supplied{TypeName: "*bytes.Buffer", Err: errors.New("great sadness")},
			want:  "[Pi] ERROR	Failed to supply *bytes.Buffer: great sadness\n",
		},
		{
			name:  "SuppliedError/rich error",
			input: &Supplied{TypeName: "*bytes.Buffer", Err: &richError{}},
			want:  "[Pi] ERROR	Failed to supply *bytes.Buffer: rich error\n",
		},
		{
			name: "ResourceAdded",
			input: &ResourceAdded{
				ConstructorName: "bytes.NewBuffer()",
				OutputTypeNames: []string{"*bytes.Buffer"},
				Private:         false,
				StackTrace:      []string{"main.main", "runtime.main"},
				ModuleTrace:     []string{"main.main"},
			},
			want: "[Pi] ADD RESOURCE	*bytes.Buffer <= bytes.NewBuffer()\n",
		},
		{
			name: "Resource Add privately",
			input: &ResourceAdded{
				ConstructorName: "bytes.NewBuffer()",
				OutputTypeNames: []string{"*bytes.Buffer"},
				Private:         true,
				StackTrace:      []string{"main.main", "runtime.main"},
				ModuleTrace:     []string{"main.main"},
			},
			want: "[Pi] ADD RESOURCE (PRIVATE)	*bytes.Buffer <= bytes.NewBuffer()\n",
		},
		{
			name: "Resource add with module",
			input: &ResourceAdded{
				ConstructorName: "bytes.NewBuffer()",
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
				Private:         false,
				StackTrace:      []string{"main.main", "runtime.main"},
				ModuleTrace:     []string{"main.main", "mypackage.GetMyModule"},
			},
			want: "[Pi] ADD RESOURCE	*bytes.Buffer <= bytes.NewBuffer() from module \"myModule\"\n",
		},
		{
			name: "Resource add with module privately",
			input: &ResourceAdded{
				ConstructorName: "bytes.NewBuffer()",
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
				Private:         true,
				StackTrace:      []string{"main.main", "runtime.main"},
				ModuleTrace:     []string{"main.main", "mypackage.GetMyModule"},
			},
			want: "[Pi] ADD RESOURCE (PRIVATE)	*bytes.Buffer <= bytes.NewBuffer() from module \"myModule\"\n",
		},
		{
			name: "Replaced",
			input: &Replaced{
				OutputTypeNames: []string{"*bytes.Buffer"},
				StackTrace:      []string{"main.main", "runtime.main"},
				ModuleTrace:     []string{"main.main"},
			},
			want: "[Pi] REPLACE	*bytes.Buffer\n",
		},
		{
			name: "Replaced with module",
			input: &Replaced{
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
				StackTrace:      []string{"main.main", "runtime.main"},
				ModuleTrace:     []string{"main.main", "mypackage.GetMyModule"},
			},
			want: "[Pi] REPLACE	*bytes.Buffer from module \"myModule\"\n",
		},
		{
			name:  "ReplacedError",
			input: &Replaced{Err: errors.New("some error")},
			want:  "[Pi] ERROR	Failed to replace: some error\n",
		},
		{
			name: "Decorated",
			input: &Decorated{
				DecoratorName:   "bytes.NewBuffer()",
				OutputTypeNames: []string{"*bytes.Buffer"},
				StackTrace:      []string{"main.main", "runtime.main"},
				ModuleTrace:     []string{"main.main"},
			},
			want: "[Pi] DECORATE	*bytes.Buffer <= bytes.NewBuffer()\n",
		},
		{
			name: "Decorated with module",
			input: &Decorated{
				DecoratorName:   "bytes.NewBuffer()",
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
				StackTrace:      []string{"main.main", "runtime.main"},
				ModuleTrace:     []string{"main.main", "mypackage.GetMyModule"},
			},
			want: "[Pi] DECORATE	*bytes.Buffer <= bytes.NewBuffer() from module \"myModule\"\n",
		},
		{
			name:  "DecorateError",
			input: &Decorated{Err: errors.New("some error")},
			want:  "[Pi] Error after options were applied: some error\n",
		},
		{
			name:  "DecorateError/rich error",
			input: &Decorated{Err: &richError{}},
			want:  "[Pi] Error after options were applied: rich error\n",
		},
		{
			name:  "Run",
			input: &Run{Name: "bytes.NewBuffer()", Kind: "constructor", ElapsedTime: 10 * time.Nanosecond},
			want:  "[Pi] RUN\tconstructor: bytes.NewBuffer() in 10ns\n",
		},
		{
			name: "Run with module",
			input: &Run{
				Name:        "bytes.NewBuffer()",
				Kind:        "constructor",
				ElapsedTime: 50 * time.Millisecond,
				ModuleName:  "myModule",
			},
			want: "[Pi] RUN\tconstructor: bytes.NewBuffer() in 50ms from module \"myModule\"\n",
		},
		{
			name: "RunError",
			input: &Run{
				Name:        "bytes.NewBuffer()",
				Kind:        "constructor",
				ElapsedTime: 5 * time.Second,
				Err:         errors.New("terrible constructor error"),
			},
			want: joinLines(
				"[Pi] RUN\tconstructor: bytes.NewBuffer() in 5s",
				"[Pi] Error returned: terrible constructor error",
			),
		},
		{
			name:  "SettingUp",
			input: &SettingUp{FunctionName: "bytes.NewBuffer()"},
			want:  "[Pi] SETTINGUP		bytes.NewBuffer()\n",
		},
		{
			name: "SettingUp with module",
			input: &SettingUp{
				FunctionName: "bytes.NewBuffer()",
				ModuleName:   "myModule",
			},
			want: "[Pi] SETTINGUP		bytes.NewBuffer() from module \"myModule\"\n",
		},
		{
			name: "SetupCompleted/Error",
			input: &SetupCompleted{
				FunctionName: "bytes.NewBuffer()",
				Err:          errors.New("some error"),
				Trace:        "foo()\n\tbar/baz.go:42\n",
			},
			want: joinLines(
				"[Pi] ERROR		pi.Setup(bytes.NewBuffer()) called from:",
				"foo()",
				"	bar/baz.go:42",
				"Failed: some error",
			),
		},
		{
			name: "SetupCompleted/Error/rich",
			input: &SetupCompleted{
				FunctionName: "bytes.NewBuffer()",
				Err:          &richError{},
				Trace:        "foo()\n\tbar/baz.go:42\n",
			},
			want: joinLines(
				"[Pi] ERROR		pi.Setup(bytes.NewBuffer()) called from:",
				"foo()",
				"	bar/baz.go:42",
				"Failed: rich error",
			),
		},
		{
			name:  "StartError",
			input: &Started{Err: errors.New("some error")},
			want:  "[Pi] ERROR		Failed to start: some error\n",
		},
		{
			name:  "StartError/rich error",
			input: &Started{Err: &richError{}},
			want:  "[Pi] ERROR		Failed to start: rich error\n",
		},
		{
			name:  "Stopping",
			input: &Stopping{Signal: os.Interrupt},
			want:  "[Pi] INTERRUPT\n",
		},
		{
			name:  "Stopped",
			input: &Stopped{Err: errors.New("some error")},
			want:  "[Pi] ERROR		Failed to stop cleanly: some error\n",
		},
		{
			name:  "Stopped/rich error",
			input: &Stopped{Err: &richError{}},
			want:  "[Pi] ERROR		Failed to stop cleanly: rich error\n",
		},
		{
			name:  "RollingBack",
			input: &RollingBack{StartErr: errors.New("some error")},
			want:  "[Pi] ERROR		Start failed, rolling back: some error\n",
		},
		{
			name:  "RollingBack/rich error",
			input: &RollingBack{StartErr: &richError{}},
			want:  "[Pi] ERROR		Start failed, rolling back: rich error\n",
		},
		{
			name:  "RolledBack",
			input: &RolledBack{Err: errors.New("some error")},
			want:  "[Pi] ERROR		Couldn't roll back cleanly: some error\n",
		},
		{
			name:  "RolledBack/rich error",
			input: &RolledBack{Err: &richError{}},
			want:  "[Pi] ERROR		Couldn't roll back cleanly: rich error\n",
		},
		{
			name:  "Started",
			input: &Started{},
			want:  "[Pi] RUNNING\n",
		},
		{
			name:  "CustomLoggerError",
			input: &LoggerInitialized{Err: errors.New("great sadness")},
			want:  "[Pi] ERROR		Failed to initialize custom logger: great sadness\n",
		},
		{
			name:  "CustomLoggerError/rich error",
			input: &LoggerInitialized{Err: &richError{}},
			want:  "[Pi] ERROR		Failed to initialize custom logger: rich error\n",
		},
		{
			name:  "LoggerInitialized",
			input: &LoggerInitialized{ConstructorName: "github.com/n24-x/pi/pievent.TestConsoleLogger.func1()"},
			want:  "[Pi] LOGGER	Initialized custom logger from github.com/n24-x/pi/pievent.TestConsoleLogger.func1()\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buff bytes.Buffer
			(&ConsoleLogger{W: &buff}).LogEvent(tt.input)

			assert.Equal(t, tt.want, buff.String())
		})
	}
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}
