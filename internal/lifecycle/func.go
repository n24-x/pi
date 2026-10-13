package lifecycle

import (
	"context"
	"reflect"

	"github.com/n24-x/pi/internal/pireflect"
)

// Reflection types for each of the supported hook function signatures. These
// are used in cases in which the Callable constraint matches a user-defined
// function type that cannot be converted to an underlying function type with
// a conventional conversion or type switch.
var (
	_reflFunc             = reflect.TypeFor[Func]()
	_reflErrorFunc        = reflect.TypeFor[ErrorFunc]()
	_reflContextFunc      = reflect.TypeFor[ContextFunc]()
	_reflContextErrorFunc = reflect.TypeFor[ContextErrorFunc]()
)

// Discrete function signatures that are allowed as part of a [Callable].
type (
	// A Func can be converted to a ContextErrorFunc.
	Func = func()
	// An ErrorFunc can be converted to a ContextErrorFunc.
	ErrorFunc = func() error
	// A ContextFunc can be converted to a ContextErrorFunc.
	ContextFunc = func(context.Context)
	// A ContextErrorFunc is used as a [Hook.OnStart] or [Hook.OnStop]
	// function.
	ContextErrorFunc = func(context.Context) error
)

// A Callable is a constraint that matches functions that are, or can be
// converted to, functions suitable for a Hook.
//
// Callable must be identical to [fx.HookFunc].
type Callable interface {
	~Func | ~ErrorFunc | ~ContextFunc | ~ContextErrorFunc
}

// Wrap wraps x into a ContextErrorFunc suitable for a Hook.
func Wrap[T Callable](x T) (ContextErrorFunc, string) {
	if x == nil {
		return nil, ""
	}

	switch fn := any(x).(type) {
	case Func:
		return func(context.Context) error {
			fn()
			return nil
		}, pireflect.FuncName(x)
	case ErrorFunc:
		return func(context.Context) error {
			return fn()
		}, pireflect.FuncName(x)
	case ContextFunc:
		return func(ctx context.Context) error {
			fn(ctx)
			return nil
		}, pireflect.FuncName(x)
	case ContextErrorFunc:
		return fn, pireflect.FuncName(x)
	}

	// Since (1) we're already using reflect in Fx, (2) we're not particularly
	// concerned with performance, and (3) unsafe would require discrete build
	// targets for appengine (etc), just use reflect to convert user-defined
	// function types to their underlying function types and then call Wrap
	// again with the converted value.
	reflVal := reflect.ValueOf(x)
	switch {
	case reflVal.CanConvert(_reflFunc):
		return Wrap(reflVal.Convert(_reflFunc).Interface().(Func))
	case reflVal.CanConvert(_reflErrorFunc):
		return Wrap(reflVal.Convert(_reflErrorFunc).Interface().(ErrorFunc))
	case reflVal.CanConvert(_reflContextFunc):
		return Wrap(reflVal.Convert(_reflContextFunc).Interface().(ContextFunc))
	default:
		// Is already convertible to ContextErrorFunc.
		return Wrap(reflVal.Convert(_reflContextErrorFunc).Interface().(ContextErrorFunc))
	}
}
