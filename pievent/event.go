package pievent

import (
	"os"
	"time"
)

// Event defines an event emitted by pi.
type Event interface {
	event() // Only pilog can implement this interface.
}

// Passing events by type to make Event hashable in the future.
func (*OnStartExecuting) event()  {}
func (*OnStartExecuted) event()   {}
func (*OnStopExecuting) event()   {}
func (*OnStopExecuted) event()    {}
func (*Supplied) event()          {}
func (*ResourceAdded) event()     {}
func (*Replaced) event()          {}
func (*Decorated) event()         {}
func (*BeforeRun) event()         {}
func (*Run) event()               {}
func (*SettingUp) event()         {}
func (*SetupCompleted) event()    {}
func (*Stopping) event()          {}
func (*Stopped) event()           {}
func (*RollingBack) event()       {}
func (*RolledBack) event()        {}
func (*Started) event()           {}
func (*LoggerInitialized) event() {}

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

// TODO It may be useless and requires review.
//
// Supplied is emitted after a value is added with pi.Supply.
type Supplied struct {
	// TypeName is the name of the type of value that was added.
	TypeName string

	// StackTrace is the stack trace of the call to Supply.
	StackTrace []string

	// ModuleTrace contains the module locations through which this value was added.
	ModuleTrace []string

	// ModuleName is the name of the module in which the value was added to.
	ModuleName string

	// Err is non-nil if we failed to supply the value.
	Err error
}

// ResourceAdded is emitted when a constructor is added to Pi.
type ResourceAdded struct {
	// ConstructorName is the name of the constructor that was provided to
	// Pi.
	ConstructorName string

	// StackTrace is the stack trace of where the constructor was provided to Pi.
	StackTrace []string

	// ModuleTrace contains the module locations through which this was provided to Pi.
	ModuleTrace []string

	// OutputTypeNames is a list of names of types that are produced by
	// this constructor.
	OutputTypeNames []string

	// ModuleName is the name of the module in which the constructor was
	// provided to.
	ModuleName string

	// Err is non-nil if we failed to provide this constructor.
	Err error

	// Private denotes whether the provided constructor is a [Private] constructor.
	Private bool
}

// Replaced is emitted when a value replaces a type in Pi.
type Replaced struct {
	// OutputTypeNames is a list of names of types that were replaced.
	OutputTypeNames []string

	// StackTrace is the stack trace of the call to Replace.
	StackTrace []string

	// ModuleTrace contains the module locations through which this value was added.
	ModuleTrace []string

	// ModuleName is the name of the module in which the value was added to.
	ModuleName string

	// Err is non-nil if we failed to supply the value.
	Err error
}

// Decorated is emitted when a decorator is executed in Pi.
type Decorated struct {
	// DecoratorName is the name of the decorator function that was
	// provided to Pi.
	DecoratorName string

	// StackTrace is the stack trace of where the decorator was given to Pi.
	StackTrace []string

	// ModuleTrace contains the module locations through which this value was added.
	ModuleTrace []string

	// ModuleName is the name of the module in which the value was added to.
	ModuleName string

	// OutputTypeNames is a list of names of types that are decorated by
	// this decorator.
	OutputTypeNames []string

	// Err is non-nil if we failed to run this decorator.
	Err error
}

// BeforeRun is emitted before a constructor, decorator, or supply/replace stub is run by Pi.
// When complete, a Run will be emitted.
type BeforeRun struct {
	// Name is the name of the function that will be run.
	Name string

	// Kind indicates which Pi option was used to pass along the function.
	// It is either "resource", "decorate", "supply", or "replace".
	Kind string

	// ModuleName is the name of the module in which the function belongs.
	ModuleName string
}

// Run is emitted after a constructor, decorator, or supply/replace stub is run by Pi.
type Run struct {
	// Name is the name of the function that was run.
	Name string

	// Kind indicates which Pi option was used to pass along the function.
	// It is either "resource", "decorate", "supply", or "replace".
	Kind string

	// ModuleName is the name of the module in which the function belongs.
	ModuleName string

	// ElapsedTime specifies how long it took to run this function.
	ElapsedTime time.Duration

	// TODO
	// Err is non-nil if the function returned an error.
	// If pi.RecoverFromPanics is used, this will include panics.
	Err error
}

// SettingUp is emitted before we invoke a function specified with pi.Setup.
type SettingUp struct {
	// FunctionName is the name of the function that will be invoked.
	FunctionName string

	// ModuleName is the name of the module in which the value was added to.
	ModuleName string
}

// SetupCompleted is emitted after we invoke a function specified with pi.Setup,
// whether it succeeded or failed.
type SetupCompleted struct {
	// Functionname is the name of the function that was invoked.
	FunctionName string

	// ModuleName is the name of the module in which the value was added to.
	ModuleName string

	// Err is non-nil if the function failed to execute.
	Err error

	// Trace records information about where the pi.Setup call was made.
	// Note that this is NOT a stack trace of the error itself.
	Trace string
}

// Started is emitted when an application is started successfully and/or it
// errored.
type Started struct {
	// Err is non-nil if the application failed to start successfully.
	Err error
}

// TODO need review
//
// Stopping is emitted when the application receives a signal to shut down
// after starting. This may happen with pi.Shutdowner or by sending a signal to
// the application on the command line.
type Stopping struct {
	// Signal is the signal that caused this shutdown.
	Signal os.Signal
}

// Stopped is emitted when the application has finished shutting down, whether
// successfully or not.
type Stopped struct {
	// Err is non-nil if errors were encountered during shutdown.
	Err error
}

// RollingBack is emitted when the application failed to start up due to an
// error, and is being rolled back.
type RollingBack struct {
	// StartErr is the error that caused this rollback.
	StartErr error
}

// RolledBack is emitted after a service has been rolled back, whether it
// succeeded or not.
type RolledBack struct {
	// Err is non-nil if the rollback failed.
	Err error
}

// LoggerInitialized is emitted when a logger supplied with pi.WithLogger is
// instantiated, or if it fails to instantiate.
type LoggerInitialized struct {
	// ConstructorName is the name of the constructor that builds this
	// logger.
	ConstructorName string

	// Err is non-nil if the logger failed to build.
	Err error
}
