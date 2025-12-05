package day04

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestDay4Solver_getLargestPossibleJoltage(t *testing.T) {
	t.Run("Correctly calculates the number of rolls where there are fewer than four rolls of paper in the eight adjacent positions", func(t *testing.T) {
		t.Parallel()
		//given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()

		day3Solver, _ := NewDay4Solver(logger)

		//when
		input := `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

		result, err := day3Solver.Solve(context.Background(), strings.NewReader(input), RemovalModeSingleLayer)

		//then
		assert.NoError(t, err)
		assert.Equal(t, 13, result)
	})

	t.Run("Correctly calculates the number of rolls where there are fewer than four rolls of paper in the eight adjacent positions and we remove reachable rolls", func(t *testing.T) {
		t.Parallel()
		//given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()

		day3Solver, _ := NewDay4Solver(logger)

		//when
		input := `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

		result, err := day3Solver.Solve(context.Background(), strings.NewReader(input), RemovalModeRecursive)

		//then
		assert.NoError(t, err)
		assert.Equal(t, 43, result)
	})
}
