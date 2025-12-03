package day03

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"
)

type Day3Solver struct {
	logger       *zap.Logger
	batteryCount int
}

func NewDay3Solver(batteryCount int, logger *zap.Logger) (*Day3Solver, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	day3Solver := &Day3Solver{
		batteryCount: batteryCount,
		logger:       logger,
	}
	return day3Solver, nil
}

func (d *Day3Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
	solution := 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		largestPossibleJoltage, err := d.getLargesPossibleJoltage(line, d.batteryCount)
		if err != nil {
			return 0, err
		}
		d.logger.Debug("got largest joltage for power bank",
			zap.Int("joltage", largestPossibleJoltage),
			zap.String("powerBank", line),
		)
		solution += largestPossibleJoltage
	}
	return solution, nil
}

// NOTE: Time complexity O(n * k) time complexity solution
// space complexiy is O(1)
func (d *Day3Solver) getLargesPossibleJoltage(powerBank string, batteryCount int) (int, error) {
	if batteryCount > len(powerBank) {
		return 0, fmt.Errorf("battery count %d is greater than the size of power bank %d", batteryCount, len(powerBank))
	}

	sum := 0
	lastPickedPos := -1

	for batteryIdx := range batteryCount {
		minPos := lastPickedPos + 1
		maxPos := len(powerBank) - batteryCount + batteryIdx

		bestDigit := -1
		bestPos := -1

		for pos := minPos; pos <= maxPos; pos++ {
			digit := int(powerBank[pos] - '0')
			if digit > bestDigit {
				bestDigit = digit
				bestPos = pos
			}
		}

		sum = sum*10 + bestDigit
		lastPickedPos = bestPos
	}

	return sum, nil
}
