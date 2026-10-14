package pilog

import (
	"fmt"
	"testing"

	"github.com/n24-x/pi/pievent"

	"github.com/stretchr/testify/assert"
)

func TestSpy(t *testing.T) {
	t.Parallel()

	var s Spy

	t.Run("empty spy", func(t *testing.T) {
		assert.Empty(t, s.Events(), "events must be empty")
		assert.Zero(t, s.Events().Len(), "events length must be zero")
		assert.Empty(t, s.EventTypes(), "event types must be empty")
	})

	s.LogEvent(&pievent.Started{})
	t.Run("use after reset", func(t *testing.T) {
		assert.Equal(t, "Started", s.EventTypes()[0])
	})

	s.LogEvent(&pievent.ResourceAdded{Err: fmt.Errorf("some error")})
	t.Run("some error", func(t *testing.T) {
		assert.Equal(t, 1, s.Events().SelectByTypeName("ResourceAdded").Len())
		assert.Equal(t, "ResourceAdded", s.EventTypes()[1])
	})

	s.Reset()
	t.Run("reset", func(t *testing.T) {
		assert.Empty(t, s.Events(), "events must be empty")
		assert.Empty(t, s.EventTypes(), "event types must be empty")
	})

	s.LogEvent(&pievent.Started{})
	t.Run("use after reset", func(t *testing.T) {
		assert.Equal(t, "Started", s.EventTypes()[0])
	})
}
