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
		largestPossibleJoltage, err := d.getLargesPossibleJoltage(line, 2)
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

func (d *Day3Solver) getLargesPossibleJoltage(powerBank string, batteryCount int) (int, error) {
	powerBankValueMap := make(map[int][]int)

	for i, powerBankCellAsRune := range powerBank {
		powerBankCell := int(powerBankCellAsRune - '0')
		val, ok := powerBankValueMap[powerBankCell]
		if !ok {
			val = make([]int, 0)
		}
		val = append(val, i)
		powerBankValueMap[powerBankCell] = val
	}

	batteries := make([]int, batteryCount)

	for batteryIdx := range batteries {
	CheckLargest:
		for i := 9; i > 0; i-- {
			powerBankValueIndexes, ok := powerBankValueMap[i]
			// if the value does not exist in the power bank skip it
			if !ok {
				continue CheckLargest
			}
			for _, powerBankValueIndex := range powerBankValueIndexes {
				spaceRequired := batteryCount - batteryIdx
				spaceAvailable := len(powerBank) - powerBankValueIndex
				// if there is not enough space for the reamining batteries we can't use this value
				if spaceRequired > spaceAvailable {
					continue CheckLargest
				}
				// if the position we want to pick comes before the position we already picked then it is invalid
				if batteryIdx > 0 && batteries[batteryIdx-1] >= powerBankValueIndex {
					continue
				}
				batteries[batteryIdx] = powerBankValueIndex
				break CheckLargest
			}
		}
	}

	sum := 0
	decimal := 1

	for i := len(batteries) - 1; i >= 0; i-- {
		powerBankValueIdx := batteries[i]
		powerBankValue := int(powerBank[powerBankValueIdx] - '0')
		sum += decimal * powerBankValue
		decimal = 10 * decimal
	}

	return sum, nil
}
