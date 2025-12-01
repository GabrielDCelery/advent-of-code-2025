package day_01

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	// given
	reader := strings.NewReader("L68\nL30\nR48\nL5\nR60\nL55\nL1\nL99\nR14\nL82")
	dial := NewDial()

	// when
	result, err := dial.GetPassword(reader)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 3, result)
}
