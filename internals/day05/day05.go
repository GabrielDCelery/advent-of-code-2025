package day05

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type ReadMode int

const (
	ReadingRanges ReadMode = iota
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

type Solution struct {
	FreshIngredientsCount     int
	AvailableIngredientsCount int
}

func (d *Day5Solver) Solve(ctx context.Context, reader io.Reader) (Solution, error) {
	var readMode = ReadingRanges

	scanner := bufio.NewScanner(reader)

	ingredientRanges := IngredientRanges{}

	solution := Solution{
		FreshIngredientsCount:     0, // count of ingredients with valid ranges
		AvailableIngredientsCount: 0, // total count of all unique ingredients in all the ranges
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readMode = ReadingIngredients
			ingredientRanges.merge()
			d.logger.Debug("merged ingredient ranges", zap.String("ranges", fmt.Sprintf("%+v", ingredientRanges)))
			continue
		}
		switch readMode {
		case ReadingRanges:
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
			if ingredientRanges.isFreshIngredient(ingredient) {
				solution.FreshIngredientsCount += 1
			}
		}
	}

	solution.AvailableIngredientsCount = ingredientRanges.countAvailableIngredients()

	return solution, nil
}

// IngredientRange represents a closed interval [min,max] of ingredient ID
type IngredientRange struct {
	min int
	max int
}

type IngredientRanges []IngredientRange

func (ir *IngredientRanges) addRange(rang IngredientRange) {
	*ir = append(*ir, rang)
}

func (ir *IngredientRanges) merge() {
	if len(*ir) == 1 {
		return
	}
	ranges := *ir
	sort.Slice(ranges, func(i int, j int) bool {
		return ranges[i].min < ranges[j].min
	})
	merged := make(IngredientRanges, 0)
	current := ranges[0]
	for i := 1; i < len(*ir); i++ {
		// merge overlapping ranges (e.g. [1-5] and [3-7] becomes [1-7])
		if ranges[i].min <= current.max {
			current.max = max(current.max, ranges[i].max)
		} else {
			merged = append(merged, current)
			current = ranges[i]
		}
	}
	merged = append(merged, current)
	*ir = merged
}

func (ir *IngredientRanges) isFreshIngredient(ingredient int) bool {
	for _, rang := range *ir {
		if ingredient >= rang.min && ingredient <= rang.max {
			return true
		}
	}
	return false
}

func (ir *IngredientRanges) countAvailableIngredients() int {
	numOfAvailableIngredients := 0
	for _, rang := range *ir {
		numOfAvailableIngredients += (rang.max - rang.min + 1)
	}
	return numOfAvailableIngredients
}
