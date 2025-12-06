package day06

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestDay6Solver(t *testing.T) {
	t.Run("Solves day 6 challenge", func(t *testing.T) {
		t.Parallel()
		//given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()
		solver, _ := NewDay6Solver(logger)
		input := `123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   + `

		//when
		solution, err := solver.Solve(context.Background(), strings.NewReader(input))

		//then
		assert.NoError(t, err)
		assert.Equal(t, 4277556, solution)
	})
}
