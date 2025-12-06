package day06

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type PuzzleInterpreter int

const (
	HumanMath PuzzleInterpreter = iota
	CephalopodMath
)

type Day6Solver struct {
	logger *zap.Logger
}

func NewDay6Solver(logger *zap.Logger) (*Day6Solver, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	solver := &Day6Solver{
		logger: logger,
	}
	return solver, nil
}

func (d *Day6Solver) Solve(ctx context.Context, reader io.Reader, puzzleInterpreter PuzzleInterpreter) (int, error) {
	problemsContainer := NewProblemsContainer(puzzleInterpreter)
	problemsContainer.parseInput(reader)
	d.logger.Debug("parsed input to container", zap.String("container", fmt.Sprintf("%+v", problemsContainer)))
	sum, err := problemsContainer.solve()
	if err != nil {
		return 0, err
	}
	return sum, nil
}

type Problem struct {
	id         int
	numsMatrix []string
	operator   string
	start      int
	end        int
}

func (p *Problem) parseNumsMatrixToNums(puzzleInterpreter PuzzleInterpreter) ([]int, error) {
	switch puzzleInterpreter {
	case HumanMath:
		nums := []int{}
		for _, numAsStr := range p.numsMatrix {
			numAsStr = strings.TrimSpace(numAsStr)
			num, err := strconv.Atoi(numAsStr)
			if err != nil {
				return nil, fmt.Errorf("invalid integer '%s'", numAsStr)
			}
			nums = append(nums, num)
		}
		return nums, nil
	case CephalopodMath:
		nums := []int{}
		for i := 0; i < (p.end - p.start); i++ {
			numAsStr := ""
			for j := 0; j < len(p.numsMatrix); j++ {
				char := p.numsMatrix[j][i]
				charAsStr := string(char)
				if charAsStr == " " {
					continue
				}
				numAsStr = numAsStr + charAsStr
			}
			num, err := strconv.Atoi(numAsStr)
			if err != nil {
				return nil, fmt.Errorf("invalid integer '%s'", numAsStr)
			}
			nums = append(nums, num)
		}
		return nums, nil
	}
	return nil, fmt.Errorf("invalid interpreter '%d'", puzzleInterpreter)
}

func (p *Problem) solve(puzzleInterpreter PuzzleInterpreter) (int, error) {
	nums, err := p.parseNumsMatrixToNums(puzzleInterpreter)
	if err != nil {
		return 0, err
	}
	if p.operator == "*" {
		sum := 1
		for _, num := range nums {
			sum *= num
		}
		return sum, nil
	}
	if p.operator == "+" {
		sum := 0
		for _, num := range nums {
			sum += num
		}
		return sum, nil
	}
	return 0, fmt.Errorf("invalid operator %s", p.operator)
}

type ProblemsContainer struct {
	puzzleInterpreter PuzzleInterpreter
	problems          []Problem
}

func NewProblemsContainer(puzzleInterpreter PuzzleInterpreter) *ProblemsContainer {
	return &ProblemsContainer{
		puzzleInterpreter: puzzleInterpreter,
		problems:          []Problem{},
	}
}

func (pc *ProblemsContainer) parseInput(reader io.Reader) {
	numLines := make([]string, 0)
	operatorLine := ""

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		isLastLine := strings.ContainsAny(line, `+*`)
		if isLastLine {
			operatorLine = line
		} else {
			numLines = append(numLines, line)
		}
	}

	start := 0

	for i, run := range operatorLine {
		char := string(run)
		isFirstChar := i == 0
		isOperator := char == "*" || char == "+"
		isLastChar := i == len(operatorLine)-1
		if isOperator && !isFirstChar {
			end := i - 1
			problem := Problem{
				id:         len(pc.problems),
				numsMatrix: []string{},
				operator:   string(operatorLine[start]),
				start:      start,
				end:        end,
			}
			pc.problems = append(pc.problems, problem)
			start = i
			continue
		}
		if isLastChar {
			end := i + 1
			problem := Problem{
				id:         len(pc.problems),
				numsMatrix: []string{},
				operator:   string(operatorLine[start]),
				start:      start,
				end:        end,
			}
			pc.problems = append(pc.problems, problem)
			start = end
			continue
		}
	}

	for i, problem := range pc.problems {
		for _, numLine := range numLines {
			pc.problems[i].numsMatrix = append(pc.problems[i].numsMatrix, numLine[problem.start:problem.end])
		}
	}
}

func (pc *ProblemsContainer) solve() (int, error) {
	sum := 0
	for _, problem := range pc.problems {
		result, err := problem.solve(pc.puzzleInterpreter)
		if err != nil {
			return 0, err
		}
		sum += result
	}
	return sum, nil
}
