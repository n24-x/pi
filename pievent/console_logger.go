package pievent

import (
	"fmt"
	"io"
	"strings"
)

// ConsoleLogger is a Pi event logger that attempts to write human-readable
// messages to the console.
//
// Use this during development.
type ConsoleLogger struct {
	W io.Writer
}

var _ Logger = (*ConsoleLogger)(nil)

func (l *ConsoleLogger) logf(msg string, args ...any) {
	fmt.Fprintf(l.W, "[Pi] "+msg+"\n", args...)
}

// LogEvent logs the given event to the console.
func (l *ConsoleLogger) LogEvent(event Event) {
	switch e := event.(type) {
	case *OnStartExecuting:
		l.logf("HOOK OnStart\t\t%s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *OnStartExecuted:
		if e.Err != nil {
			l.logf("HOOK OnStart\t\t%s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.ElapsedTime, e.Err)
		} else {
			l.logf("HOOK OnStart\t\t%s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.ElapsedTime)
		}
	case *OnStopExecuting:
		l.logf("HOOK OnStop\t\t%s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *OnStopExecuted:
		if e.Err != nil {
			l.logf("HOOK OnStop\t\t%s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.ElapsedTime, e.Err)
		} else {
			l.logf("HOOK OnStop\t\t%s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.ElapsedTime)
		}
	case *Supplied:
		if e.Err != nil {
			l.logf("ERROR\tFailed to supply %v: %+v", e.TypeName, e.Err)
		} else if e.ModuleName != "" {
			l.logf("SUPPLY\t%v from module %q", e.TypeName, e.ModuleName)
		} else {
			l.logf("SUPPLY\t%v", e.TypeName)
		}
	case *ResourceAdded:
		var privateStr string
		if e.Private {
			privateStr = " (PRIVATE)"
		}
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.logf("ADD RESOURCE%v\t%v <= %v from module %q", privateStr, rtype, e.ConstructorName, e.ModuleName)
			} else {
				l.logf("ADD RESOURCE%v\t%v <= %v", privateStr, rtype, e.ConstructorName)
			}
		}
		if e.Err != nil {
			l.logf("Error after options were applied: %+v", e.Err)
		}
	case *Replaced:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.logf("REPLACE\t%v from module %q", rtype, e.ModuleName)
			} else {
				l.logf("REPLACE\t%v", rtype)
			}
		}
		if e.Err != nil {
			l.logf("ERROR\tFailed to replace: %+v", e.Err)
		}
	case *Decorated:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.logf("DECORATE\t%v <= %v from module %q", rtype, e.DecoratorName, e.ModuleName)
			} else {
				l.logf("DECORATE\t%v <= %v", rtype, e.DecoratorName)
			}
		}
		if e.Err != nil {
			l.logf("Error after options were applied: %+v", e.Err)
		}
	case *BeforeRun:
		var moduleStr string
		if e.ModuleName != "" {
			moduleStr = fmt.Sprintf(" from module %q", e.ModuleName)
		}
		l.logf("BEFORE RUN\t%s: %s%s", e.Kind, e.Name, moduleStr)
	case *Run:
		var moduleStr string
		if e.ModuleName != "" {
			moduleStr = fmt.Sprintf(" from module %q", e.ModuleName)
		}
		l.logf("RUN\t%v: %v in %s%v", e.Kind, e.Name, e.ElapsedTime, moduleStr)
		if e.Err != nil {
			l.logf("Error returned: %+v", e.Err)
		}
	case *SettingUp:
		if e.ModuleName != "" {
			l.logf("SETTINGUP\t\t%s from module %q", e.FunctionName, e.ModuleName)
		} else {
			l.logf("SETTINGUP\t\t%s", e.FunctionName)
		}
	case *SetupCompleted:
		if e.Err != nil {
			l.logf("ERROR\t\tpi.Setup(%v) called from:\n%+vFailed: %+v", e.FunctionName, e.Trace, e.Err)
		}
	case *Stopping:
		l.logf("%v", strings.ToUpper(e.Signal.String()))
	case *Stopped:
		if e.Err != nil {
			l.logf("ERROR\t\tFailed to stop cleanly: %+v", e.Err)
		}
	case *RollingBack:
		l.logf("ERROR\t\tStart failed, rolling back: %+v", e.StartErr)
	case *RolledBack:
		if e.Err != nil {
			l.logf("ERROR\t\tCouldn't roll back cleanly: %+v", e.Err)
		}
	case *Started:
		if e.Err != nil {
			l.logf("ERROR\t\tFailed to start: %+v", e.Err)
		} else {
			l.logf("RUNNING")
		}
	case *LoggerInitialized:
		if e.Err != nil {
			l.logf("ERROR\t\tFailed to initialize custom logger: %+v", e.Err)
		} else {
			l.logf("LOGGER\tInitialized custom logger from %v", e.ConstructorName)
		}
	}
}
