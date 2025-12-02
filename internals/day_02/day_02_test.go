package day_02

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestDay2Solver(t *testing.T) {
	t.Run("correctly calculates the number of invalid items", func(t *testing.T) {
		t.Parallel()
		//given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()
		day2Solver := NewDay2Solver(logger)

		//when
		result, err := day2Solver.Solve(strings.NewReader("11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124"))

		//then
		assert.NoError(t, err)
		assert.Equal(t, 1227775554, result)
	})
}
