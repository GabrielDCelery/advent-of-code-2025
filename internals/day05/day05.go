package day05

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type ReadMode int

const (
	ReadigFreshRanges ReadMode = iota
	ReadingIngredients
)

type freshIngredientRange struct {
	min int
	max int
}

type Day5Solver struct {
	logger *zap.Logger
}

func NewDay5Solver(logger *zap.Logger) (*Day5Solver, error) {
	solver := &Day5Solver{
		logger: logger,
	}
	return solver, nil
}

func (d *Day5Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
	var readMode = ReadigFreshRanges

	scanner := bufio.NewScanner(reader)

	freshIngredientRanges := make([]freshIngredientRange, 0)

	numOfFreshIngredients := 0

Outer:
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readMode = ReadingIngredients
			continue
		}
		switch readMode {
		case ReadigFreshRanges:
			boundaries := strings.Split(line, "-")
			if len(boundaries) != 2 {
				return 0, fmt.Errorf("ingredient range '%s' should have min max boundaries", line)
			}
			min, err := strconv.Atoi(boundaries[0])
			if err != nil {
				return 0, fmt.Errorf("ingredient range '%s' should contain valid integers", line)
			}

			max, err := strconv.Atoi(boundaries[1])
			if err != nil {
				return 0, fmt.Errorf("ingredient range '%s' should contain valid integers", line)
			}
			freshIngredientRanges = append(freshIngredientRanges, freshIngredientRange{min, max})
		case ReadingIngredients:
			ingredient, err := strconv.Atoi(line)
			if err != nil {
				return 0, fmt.Errorf("ingredient '%s' should be a valid integer", line)
			}
			for _, freshIngredientRange := range freshIngredientRanges {
				if ingredient >= freshIngredientRange.min && ingredient <= freshIngredientRange.max {
					numOfFreshIngredients += 1
					continue Outer
				}
			}
		}
	}

	return numOfFreshIngredients, nil
}
