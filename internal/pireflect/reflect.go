package pireflect

import (
	"fmt"
	"reflect"
	"runtime"
)

// FuncName returns a funcs name
func FuncName(fn any) string {
	fnV := reflect.ValueOf(fn)
	if fnV.Kind() != reflect.Func {
		return fmt.Sprint(fn)
	}

	function := runtime.FuncForPC(fnV.Pointer())
	if function == nil {
		return fmt.Sprint(fn)
	}

	return fmt.Sprintf("%s()", sanitizeFuncName(function.Name()))
}

// Caller returns the formatted calling func name
func Caller() string {
	return CallerStack(1, 0).CallerName()
}
