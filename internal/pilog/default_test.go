package pilog

import (
	"testing"

	"github.com/n24-x/pi/internal/testutil"

	"github.com/n24-x/goleak"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		DefaultLogger(testutil.WriteSyncer{T: t})
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
