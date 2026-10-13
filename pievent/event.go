package pievent

import "time"

// Event defines an event emitted by pi.
type Event interface {
	event() // Only pilog can implement this interface.
}

// Passing events by type to make Event hashable in the future.
func (*OnStartExecuting) event() {}
func (*OnStartExecuted) event()  {}
func (*OnStopExecuting) event()  {}
func (*OnStopExecuted) event()   {}

// OnStartExecuting is emitted before an OnStart hook is executed.
type OnStartExecuting struct {
	// FunctionName is the name of the function that will be executed.
	FunctionName string

	// CallerName is the name of the function that scheduled the hook for
	// execution.
	CallerName string
}

// OnStartExecuted is emitted after an OnStart hook has been executed.
type OnStartExecuted struct {
	// FunctionName is the name of the function that was executed.
	FunctionName string

	// CallerName is the name of the function that scheduled the hook for
	// execution.
	CallerName string

	// Method specifies the kind of the hook. This is one of "OnStart" and
	// "OnStop".
	Method string

	// ElapsedTime specifies how long it took to run this hook.
	ElapsedTime time.Duration

	// Err is non-nil if the hook failed to execute.
	Err error
}

// OnStopExecuting is emitted before an OnStop hook is executed.
type OnStopExecuting struct {
	// FunctionName is the name of the function that will be executed.
	FunctionName string

	// CallerName is the name of the function that scheduled the hook for
	// execution.
	CallerName string
}

// OnStopExecuted is emitted after an OnStop hook has been executed.
type OnStopExecuted struct {
	// FunctionName is the name of the function that was executed.
	FunctionName string

	// CallerName is the name of the function that scheduled the hook for
	// execution.
	CallerName string

	// ElapsedTime specifies how long it took to run this hook.
	ElapsedTime time.Duration

	// Err is non-nil if the hook failed to execute.
	Err error
}
