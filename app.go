package pi

// An App is a modular application built around dependency injection. Most
// users will only need to use the New constructor and the all-in-one Run
// convenience method. In more unusual cases, users may need to use the Err,
// Start, Done, and Stop methods by hand instead of relying on Run.
//
// [New] creates and initializes an App. All applications begin with a
// constructor for the Lifecycle type already registered.
//
// In addition to that built-in functionality, users typically pass a handful
// of [Provide] options and one or more [Invoke] options. The Provide options
// teach the application how to instantiate a variety of types, and the Invoke
// options describe how to initialize the application.
//
// When created, the application immediately executes all the functions passed
// via Invoke options. To supply these functions with the parameters they
// need, the application looks for constructors that return the appropriate
// types; if constructors for any required types are missing or any
// invocations return an error, the application will fail to start (and Err
// will return a descriptive error message).
//
// Once all the invocations (and any required constructors) have been called,
// New returns and the application is ready to be started using Run or Start.
// On startup, it executes any OnStart hooks registered with its Lifecycle.
// OnStart hooks are executed one at a time, in order, and must all complete
// within a configurable deadline (by default, 15 seconds). For details on the
// order in which OnStart hooks are executed, see the documentation for the
// Start method.
//
// At this point, the application has successfully started up. If started via
// Run, it will continue operating until it receives a shutdown signal from
// Done (see the [App.Done] documentation for details); if started explicitly via
// Start, it will operate until the user calls Stop. On shutdown, OnStop hooks
// execute one at a time, in reverse order, and must all complete within a
// configurable deadline (again, 15 seconds by default).
// type App struct {
// 	err       error
// 	clock     fxclock.Clock
// 	lifecycle *lifecycleWrapper

// 	container *dig.Container
// 	root      *module

// 	// Timeouts used
// 	startTimeout time.Duration
// 	stopTimeout  time.Duration
// 	// Decides how we react to errors when building the graph.
// 	errorHooks []ErrorHandler
// 	validate   bool
// 	// Whether to recover from panics in Dig container
// 	recoverFromPanics bool

// 	// Used to signal shutdowns.
// 	receivers signalReceivers

// 	osExit func(code int) // os.Exit override; used for testing only
// }
