package testutil

import (
	"go.uber.org/zap/zapcore"
)

// TestingT is a subset of the testing.T interface.
type TestingT interface {
	Logf(string, ...any)
}

// WriteSyncer is a zapcore.WriteSyncer that writes to the provided test
// logger.
type WriteSyncer struct{ T TestingT }

var _ zapcore.WriteSyncer = WriteSyncer{}

// Write writes the provided bytes to the underlying TestingT.
func (w WriteSyncer) Write(bs []byte) (int, error) {
	w.T.Logf("%s", bs)
	return len(bs), nil
}

// Sync is a no-op.
func (WriteSyncer) Sync() error { return nil }
