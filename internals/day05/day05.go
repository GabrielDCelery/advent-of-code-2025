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
	NumOfFreshIngredients     int
	NumOfAvailableIngredients int
}

func (d *Day5Solver) Solve(ctx context.Context, reader io.Reader) (Solution, error) {
	var readMode = ReadigFreshRanges

	scanner := bufio.NewScanner(reader)

	ingredientRanges := IngredientRanges{}

	solution := Solution{
		NumOfFreshIngredients:     0,
		NumOfAvailableIngredients: 0,
	}

Scanner:
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readMode = ReadingIngredients
			ingredientRanges.merge()
			d.logger.Debug("merged ingredient ranges", zap.String("ranges", fmt.Sprintf("%+v", ingredientRanges)))
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
			if ingredientRanges.isFreshIngredient(ingredient) {
				solution.NumOfFreshIngredients += 1
				continue Scanner
			}
		}
	}

	solution.NumOfAvailableIngredients = ingredientRanges.countNumOfAvailableIngredients()

	return solution, nil
}

type ReadMode int

const (
	ReadigFreshRanges ReadMode = iota
	ReadingIngredients
)

type IngredientRange struct {
	min int
	max int
}

type IngredientRanges []IngredientRange

func (f *IngredientRanges) addRange(rang IngredientRange) {
	*f = append(*f, rang)
}

func (f *IngredientRanges) merge() {
	if len(*f) == 1 {
		return
	}
	ranges := *f
	sort.Slice(ranges, func(i int, j int) bool {
		return ranges[i].min < ranges[j].min
	})
	merged := make(IngredientRanges, 0)
	current := ranges[0]
	for i := 1; i < len(*f); i++ {
		if ranges[i].min <= current.max {
			current.max = max(current.max, ranges[i].max)
		} else {
			merged = append(merged, current)
			current = ranges[i]
		}
	}
	merged = append(merged, current)
	*f = merged
}

func (f *IngredientRanges) isFreshIngredient(ingredient int) bool {
	for _, rang := range *f {
		if ingredient >= rang.min && ingredient <= rang.max {
			return true
		}
	}
	return false
}

func (f *IngredientRanges) countNumOfAvailableIngredients() int {
	numOfAvailableIngredients := 0
	for _, rang := range *f {
		numOfAvailableIngredients += (rang.max - rang.min + 1)
	}
	return numOfAvailableIngredients
}
