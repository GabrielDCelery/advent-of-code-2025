package day03

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestDay3Solver_getLargestPossibleJoltage(t *testing.T) {
	testCases := []struct {
		poweBank string
		expected int
	}{
		{poweBank: "987654321111111", expected: 98},
		{poweBank: "811111111111119", expected: 89},
		{poweBank: "234234234234278", expected: 78},
		{poweBank: "818181911112111", expected: 92},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("Correctly extracts possible joltage %d from power bank %s", tt.expected, tt.poweBank), func(t *testing.T) {
			t.Parallel()
			//given
			logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
			defer logger.Sync()

			day3Solver, _ := NewDay3Solver(logger)

			//when
			result, err := day3Solver.getLargesPossibleJoltage(tt.poweBank)

			//then
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}

}

func TestDay3Solver_Solve(t *testing.T) {
	t.Run("Correctly calculates max possible joltage", func(t *testing.T) {
		t.Parallel()
		//given
		logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
		defer logger.Sync()
		day3Solver, _ := NewDay3Solver(logger)

		//when
		result, err := day3Solver.Solve(context.Background(), strings.NewReader("987654321111111\n811111111111119\n234234234234278\n818181911112111"))

		//then
		assert.NoError(t, err)
		assert.Equal(t, 357, result)
	})
}
