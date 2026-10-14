package pilog

import (
	"reflect"
	"sync"

	"github.com/n24-x/pi/pievent"
)

// Events is a list of events captured by pilog.Spy.
type Events []pievent.Event

// Len returns the number of events in this list.
func (es Events) Len() int {
	return len(es)
}

// SelectByTypeName returns a new list with only events matching the specified
// type.
func (es Events) SelectByTypeName(name string) Events {
	var out Events
	for _, e := range es {
		if reflect.TypeOf(e).Elem().Name() == name {
			out = append(out, e)
		}
	}
	return out
}

// Spy is a Pi event logger that captures emitted events and/or logged
// statements. It may be used in tests of Pi logs.
type Spy struct {
	mu     sync.RWMutex
	events Events
}

var _ pievent.Logger = &Spy{}

// LogEvent appends an Event.
func (s *Spy) LogEvent(event pievent.Event) {
	s.mu.Lock()
	s.events = append(s.events, event)
	s.mu.Unlock()
}

// Events returns all captured events.
func (s *Spy) Events() Events {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make(Events, len(s.events))
	copy(events, s.events)
	return events
}

// EventTypes returns all captured event types.
func (s *Spy) EventTypes() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	types := make([]string, len(s.events))
	for i, e := range s.events {
		types[i] = reflect.TypeOf(e).Elem().Name()
	}
	return types
}

// Reset clears all messages and events from the Spy.
func (s *Spy) Reset() {
	s.mu.Lock()
	s.events = s.events[:0]
	s.mu.Unlock()
}
