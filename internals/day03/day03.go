package day03

import (
	"bufio"
	"context"
	"io"

	"go.uber.org/zap"
)

type Day3Solver struct {
	logger *zap.Logger
}

func NewDay3Solver(logger *zap.Logger) (*Day3Solver, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	day3Solver := &Day3Solver{
		logger,
	}
	return day3Solver, nil
}

func (d *Day3Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
	solution := 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		largestPossibleJoltage, err := d.getLargesPossibleJoltage(line)
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

func (d *Day3Solver) getLargesPossibleJoltage(powerBank string) (int, error) {
	d1 := 0
	d2 := 1
	for i, char := range powerBank {
		digit := int(char - '0')
		if (i+1) != len(powerBank) && digit > d1 {
			d1 = digit
			d2 = int(powerBank[i+1] - '0')
			continue
		}
		if i != 0 && digit > d2 {
			d2 = digit
			continue
		}
	}
	return d1*10 + d2, nil
}
