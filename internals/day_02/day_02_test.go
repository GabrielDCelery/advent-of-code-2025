package day_02

import (
	"context"
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

		day2Solver, errSolver := NewDay2Solver(logger, SomeSequenceRepeatedTwice)

		ctx := context.Background()
		reader := strings.NewReader("11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124")

		//when
		result, err := day2Solver.Solve(ctx, reader)

		//then
		assert.NoError(t, errSolver)
		assert.NoError(t, err)
		assert.Equal(t, 1227775554, result)
	})

	// t.Run("correctly counts the number of invalid product IDs in a given range", func(t *testing.T) {
	// 	t.Parallel()
	// 	tests := []struct {
	// 		min      int
	// 		max      int
	// 		expected int
	// 	}{
	// 		{min: 11, max: 22, expected: 2},
	// 		{min: 95, max: 115, expected: 1},
	// 		{min: 998, max: 1012, expected: 1},
	// 		{min: 1188511880, max: 1188511890, expected: 1},
	// 		{min: 222220, max: 222224, expected: 1},
	// 		{min: 1698522, max: 1698528, expected: 0},
	// 		{min: 446443, max: 446449, expected: 1},
	// 		{min: 38593856, max: 38593862, expected: 1},
	// 	}
	//
	// 	for _, tt := range tests {
	// 		t.Run(fmt.Sprintf("correctly calculates invalid product ID count between %d and %d", tt.min, tt.max), func(t *testing.T) {
	// 			t.Parallel()
	// 			//when
	// 			count := getInvalidIDs(context.Background(), tt.min, tt.max)
	//
	// 			//then
	// 			assert.Equal(t, tt.expected, count)
	// 		})
	// 	}
	// })
}
