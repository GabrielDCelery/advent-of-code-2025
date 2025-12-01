package day_01

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	t.Run("Correctly handles input with 'end' password method", func(t *testing.T) {
		t.Parallel()
		// given
		reader := strings.NewReader("L68\nL30\nR48\nL5\nR60\nL55\nL1\nL99\nR14\nL82")
		dial, dialErr := NewDial("end", nil)

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
		reader := strings.NewReader("L68\nL30\nR48\nL5\nR60\nL55\nL1\nL99\nR14\nL82")
		dial, dialErr := NewDial("click", nil)

		// when
		result, err := dial.GetPassword(reader)

		// then
		assert.NoError(t, dialErr)
		assert.NoError(t, err)
		assert.Equal(t, 6, result)
	})
}
