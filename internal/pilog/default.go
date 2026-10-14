package pilog

import (
	"io"

	"github.com/n24-x/pi/pievent"
)

// DefaultLogger constructs a Logger out of io.Writer.
func DefaultLogger(w io.Writer) pievent.Logger {
	return &pievent.ConsoleLogger{W: w}
}
