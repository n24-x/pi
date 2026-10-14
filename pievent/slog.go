package pievent

import (
	"context"
	"log/slog"
	"strings"
)

// SlogLogger is a Pi event logger that logs events using a slog logger.
type SlogLogger struct {
	Logger *slog.Logger

	ctx        context.Context
	logLevel   slog.Level
	errorLevel *slog.Level
}

var _ Logger = (*SlogLogger)(nil)

// UseContext sets the context that will be used when logging to the slog logger.
func (l *SlogLogger) UseContext(ctx context.Context) {
	l.ctx = ctx
}

// UseLogLevel sets the level of non-error logs emitted by Pi to level.
func (l *SlogLogger) UseLogLevel(level slog.Level) {
	l.logLevel = level
}

// UseErrorLevel sets the level of error logs emitted by Pi to level.
func (l *SlogLogger) UseErrorLevel(level slog.Level) {
	l.errorLevel = &level
}

func (l *SlogLogger) filter(fields []any) []any {
	filtered := make([]any, 0, len(fields))

	for _, field := range fields {
		if field, ok := field.(slog.Attr); ok {
			if _, ok := field.Value.Any().(slogFieldSkip); ok {
				continue
			}
		}

		filtered = append(filtered, field)
	}

	return filtered
}

func (l *SlogLogger) logEvent(msg string, fields ...any) {
	l.Logger.Log(l.ctx, l.logLevel, msg, l.filter(fields)...)
}

func (l *SlogLogger) logError(msg string, fields ...any) {
	lvl := slog.LevelError
	if l.errorLevel != nil {
		lvl = *l.errorLevel
	}

	l.Logger.Log(l.ctx, lvl, msg, l.filter(fields)...)
}

// LogEvent logs the given event to the slog logger.
func (l *SlogLogger) LogEvent(event Event) {
	switch e := event.(type) {
	case *OnStartExecuting:
		l.logEvent("OnStart hook executing",
			slog.String("callee", e.FunctionName),
			slog.String("caller", e.CallerName),
		)
	case *OnStartExecuted:
		if e.Err != nil {
			l.logError("OnStart hook failed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slogErr(e.Err),
			)
		} else {
			l.logEvent("OnStart hook executed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.String("elapsed", e.ElapsedTime.String()),
			)
		}
	case *OnStopExecuting:
		l.logEvent("OnStop hook executing",
			slog.String("callee", e.FunctionName),
			slog.String("caller", e.CallerName),
		)
	case *OnStopExecuted:
		if e.Err != nil {
			l.logError("OnStop hook failed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slogErr(e.Err),
			)
		} else {
			l.logEvent("OnStop hook executed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.String("elapsed", e.ElapsedTime.String()),
			)
		}
	case *Supplied:
		if e.Err != nil {
			l.logError("error encountered while applying options",
				slog.String("type", e.TypeName),
				slogStrings("moduletrace", e.ModuleTrace),
				slogStrings("stacktrace", e.StackTrace),
				slogMaybeModuleField(e.ModuleName),
				slogErr(e.Err))
		} else {
			l.logEvent("supplied",
				slog.String("type", e.TypeName),
				slogStrings("stacktrace", e.StackTrace),
				slogStrings("moduletrace", e.ModuleTrace),
				slogMaybeModuleField(e.ModuleName),
			)
		}
	case *ResourceAdded:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("resource_added",
				slog.String("constructor", e.ConstructorName),
				slogStrings("stacktrace", e.StackTrace),
				slogStrings("moduletrace", e.ModuleTrace),
				slogMaybeModuleField(e.ModuleName),
				slog.String("type", rtype),
				slogMaybeBool("private", e.Private),
			)
		}
		if e.Err != nil {
			l.logError("error encountered while applying options",
				slogMaybeModuleField(e.ModuleName),
				slogStrings("stacktrace", e.StackTrace),
				slogStrings("moduletrace", e.ModuleTrace),
				slogErr(e.Err))
		}
	case *Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("replaced",
				slogStrings("stacktrace", e.StackTrace),
				slogStrings("moduletrace", e.ModuleTrace),
				slogMaybeModuleField(e.ModuleName),
				slog.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.logError("error encountered while replacing",
				slogStrings("stacktrace", e.StackTrace),
				slogStrings("moduletrace", e.ModuleTrace),
				slogMaybeModuleField(e.ModuleName),
				slogErr(e.Err))
		}
	case *Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("decorated",
				slog.String("decorator", e.DecoratorName),
				slogStrings("stacktrace", e.StackTrace),
				slogStrings("moduletrace", e.ModuleTrace),
				slogMaybeModuleField(e.ModuleName),
				slog.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.logError("error encountered while applying options",
				slogStrings("stacktrace", e.StackTrace),
				slogStrings("moduletrace", e.ModuleTrace),
				slogMaybeModuleField(e.ModuleName),
				slogErr(e.Err))
		}
	case *BeforeRun:
		l.logEvent("before run",
			slog.String("name", e.Name),
			slog.String("kind", e.Kind),
			slogMaybeModuleField(e.ModuleName),
		)
	case *Run:
		if e.Err != nil {
			l.logError("error returned",
				slog.String("name", e.Name),
				slog.String("kind", e.Kind),
				slogMaybeModuleField(e.ModuleName),
				slogErr(e.Err),
			)
		} else {
			l.logEvent("run",
				slog.String("name", e.Name),
				slog.String("kind", e.Kind),
				slog.String("elapsed", e.ElapsedTime.String()),
				slogMaybeModuleField(e.ModuleName),
			)
		}
	case *SettingUp:
		// Do not log stack as it will make logs hard to read.
		l.logEvent("settingup",
			slog.String("function", e.FunctionName),
			slogMaybeModuleField(e.ModuleName),
		)
	case *SetupCompleted:
		if e.Err != nil {
			l.logError("setup failed",
				slogErr(e.Err),
				slog.String("stack", e.Trace),
				slog.String("function", e.FunctionName),
				slogMaybeModuleField(e.ModuleName),
			)
		}
	case *Stopping:
		l.logEvent("received signal",
			slog.String("signal", strings.ToUpper(e.Signal.String())))
	case *Stopped:
		if e.Err != nil {
			l.logError("stop failed", slogErr(e.Err))
		}
	case *RollingBack:
		l.logError("start failed, rolling back", slogErr(e.StartErr))
	case *RolledBack:
		if e.Err != nil {
			l.logError("rollback failed", slogErr(e.Err))
		}
	case *Started:
		if e.Err != nil {
			l.logError("start failed", slogErr(e.Err))
		} else {
			l.logEvent("started")
		}
	case *LoggerInitialized:
		if e.Err != nil {
			l.logError("custom logger initialization failed", slogErr(e.Err))
		} else {
			l.logEvent("initialized custom pievent.Logger", slog.String("function", e.ConstructorName))
		}
	}
}
