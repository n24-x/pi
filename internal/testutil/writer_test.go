package testutil

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// spy is a TestingT that captures Log requests.
type spy struct{ Logs []string }

func (s *spy) Clear() { s.Logs = nil }

func (s *spy) Logf(msg string, args ...any) {
	s.Logs = append(s.Logs, fmt.Sprintf(msg, args...))
}

func TestWriteSyncer(t *testing.T) {
	t.Parallel()

	var spy spy
	ws := WriteSyncer{T: &spy}

	t.Run("log", func(t *testing.T) {
		t.Parallel()

		spy.Clear()

		io.WriteString(ws, "hello")
		assert.Equal(t, []string{"hello"}, spy.Logs)
	})

	t.Run("sync", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, ws.Sync())
	})
}
