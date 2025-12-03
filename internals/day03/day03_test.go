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
	t.Parallel()

	testCases := []struct {
		poweBank     string
		batteryCount int
		expected     int
	}{
		{poweBank: "987654321111111", batteryCount: 2, expected: 98},
		{poweBank: "811111111111119", batteryCount: 2, expected: 89},
		{poweBank: "234234234234278", batteryCount: 2, expected: 78},
		{poweBank: "818181911112111", batteryCount: 2, expected: 92},
		{poweBank: "987654321111111", batteryCount: 12, expected: 987654321111},
		{poweBank: "811111111111119", batteryCount: 12, expected: 811111111119},
		{poweBank: "234234234234278", batteryCount: 12, expected: 434234234278},
		{poweBank: "818181911112111", batteryCount: 12, expected: 888911112111},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("Correctly extracts possible joltage %d from power bank %s with battery count %d", tt.expected, tt.poweBank, tt.batteryCount), func(t *testing.T) {
			t.Parallel()
			//given
			logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
			defer logger.Sync()

			day3Solver, _ := NewDay3Solver(tt.batteryCount, logger)

			//when
			result, err := day3Solver.getLargesPossibleJoltage(tt.poweBank, tt.batteryCount)

			//then
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}

}

func TestDay3Solver_Solve(t *testing.T) {
	t.Run("Correctly calculates max possible joltage", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			input        string
			batteryCount int
			expected     int
		}{
			{input: "987654321111111\n811111111111119\n234234234234278\n818181911112111", batteryCount: 2, expected: 357},
			{input: "987654321111111\n811111111111119\n234234234234278\n818181911112111", batteryCount: 12, expected: 3121910778619},
		}

		for _, tt := range testCases {
			t.Run(fmt.Sprintf("returns %d for battery count %d", tt.expected, tt.batteryCount), func(t *testing.T) {
				//given
				logger := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
				defer logger.Sync()
				day3Solver, _ := NewDay3Solver(tt.batteryCount, logger)

				//when
				result, err := day3Solver.Solve(context.Background(), strings.NewReader(tt.input))

				//then
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			})
		}

	})
}
