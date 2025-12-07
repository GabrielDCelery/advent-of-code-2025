package day07

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestDay7Solver(t *testing.T) {
	t.Run("Solves day 6 challenge part 1", func(t *testing.T) {
		t.Parallel()
		//given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()
		solver, _ := NewDay7Solver(logger)
		input := `.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............`

		//when
		solution, err := solver.Solve(context.Background(), strings.NewReader(input))

		//then
		assert.NoError(t, err)
		assert.Equal(t, 21, solution)
	})
}
