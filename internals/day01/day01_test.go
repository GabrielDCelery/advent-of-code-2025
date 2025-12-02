package day01

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestDial(t *testing.T) {
	t.Run("Correctly handles input with 'end' password method", func(t *testing.T) {
		t.Parallel()
		// given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()
		reader := strings.NewReader("L68\nL30\nR48\nL5\nR60\nL55\nL1\nL99\nR14\nL82")
		dial, dialErr := NewDial("end", logger)

		// when
		result, err := dial.GetPassword(reader)

		// then
		assert.NoError(t, dialErr)
		assert.NoError(t, err)
		assert.Equal(t, 3, result)
	})

	t.Run("Correctly handles input with 'click' password method", func(t *testing.T) {
		t.Parallel()
		// given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()
		reader := strings.NewReader("L68\nL30\nR48\nL5\nR60\nL55\nL1\nL99\nR14\nL82")
		dial, dialErr := NewDial("click", logger)

		// when
		result, err := dial.GetPassword(reader)

		// then
		assert.NoError(t, dialErr)
		assert.NoError(t, err)
		assert.Equal(t, 6, result)
	})
}
