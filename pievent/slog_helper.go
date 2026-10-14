package pievent

import (
	"log/slog"
	"strconv"
)

type slogFieldSkip struct{}

func slogMaybeModuleField(name string) slog.Attr {
	if len(name) == 0 {
		return slog.Any("module", slogFieldSkip{})
	}
	return slog.String("module", name)
}

func slogMaybeBool(name string, b bool) slog.Attr {
	if !b {
		return slog.Any(name, slogFieldSkip{})
	}
	return slog.Bool(name, true)
}

func slogErr(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func slogStrings(key string, str []string) slog.Attr {
	attrs := make([]any, len(str))
	for i, val := range str {
		attrs[i] = slog.String(strconv.Itoa(i), val)
	}
	return slog.Group(key, attrs...)
}
