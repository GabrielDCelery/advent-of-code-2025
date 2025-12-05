package day05

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestDay5Solver(t *testing.T) {
	t.Run("Correctly calculates the number of fresh ingredients", func(t *testing.T) {
		//given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()
		solver, _ := NewDay5Solver(logger)
		input := `3-5
10-14
16-20
12-18

1
5
8
11
17
32`

		//when
		result, err := solver.Solve(context.Background(), strings.NewReader(input))

		//then
		assert.NoError(t, err)
		assert.Equal(t, 3, result.NumOfFreshIngredients)
	})

	t.Run("Correctly calculates the number of available ingredients", func(t *testing.T) {
		//given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()
		solver, _ := NewDay5Solver(logger)
		input := `3-5
10-14
16-20
12-18

1
5
8
11
17
32`

		//when
		result, err := solver.Solve(context.Background(), strings.NewReader(input))

		//then
		assert.NoError(t, err)
		assert.Equal(t, 14, result.NumOfAvailableIngredients)
	})
}
