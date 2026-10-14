package pievent

// Logger defines interface used for logging.
type Logger interface {
	// LogEvent is called when a logging event is emitted.
	LogEvent(Event)
}

// NopLogger is a Pi event logger that ignores all messages.
var NopLogger = nopLogger{}

type nopLogger struct{}

var _ Logger = nopLogger{}

func (nopLogger) LogEvent(Event) {}

func (nopLogger) String() string { return "NopLogger" }
