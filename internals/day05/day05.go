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

type Day5Solver struct {
	logger *zap.Logger
}

func NewDay5Solver(logger *zap.Logger) (*Day5Solver, error) {
	solver := &Day5Solver{
		logger: logger,
	}
	return solver, nil
}

type IngredientRange struct {
	min int
	max int
}

type IngredientRanges []IngredientRange

func (f *IngredientRanges) addRange(rang IngredientRange) {
	*f = append(*f, rang)
}

func (f *IngredientRanges) merge() {
	*f = *f
}

type Solution struct {
	NumOfFreshIngredients int
}

func (d *Day5Solver) Solve(ctx context.Context, reader io.Reader) (Solution, error) {
	var readMode = ReadigFreshRanges

	scanner := bufio.NewScanner(reader)

	ingredientRanges := IngredientRanges{}

	solution := Solution{
		NumOfFreshIngredients: 0,
	}

Outer:
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readMode = ReadingIngredients
			ingredientRanges.merge()
			continue
		}
		switch readMode {
		case ReadigFreshRanges:
			boundaries := strings.Split(line, "-")
			if len(boundaries) != 2 {
				return Solution{}, fmt.Errorf("ingredient range '%s' should have min max boundaries", line)
			}
			min, err := strconv.Atoi(boundaries[0])
			if err != nil {
				return Solution{}, fmt.Errorf("ingredient range '%s' should contain valid integers", line)
			}

			max, err := strconv.Atoi(boundaries[1])
			if err != nil {
				return Solution{}, fmt.Errorf("ingredient range '%s' should contain valid integers", line)
			}
			ingredientRanges.addRange(IngredientRange{min, max})
		case ReadingIngredients:
			ingredient, err := strconv.Atoi(line)
			if err != nil {
				return Solution{}, fmt.Errorf("ingredient '%s' should be a valid integer", line)
			}
			for _, ingredientRange := range ingredientRanges {
				if ingredient >= ingredientRange.min && ingredient <= ingredientRange.max {
					solution.NumOfFreshIngredients += 1
					continue Outer
				}
			}
		}
	}

	return solution, nil
}
