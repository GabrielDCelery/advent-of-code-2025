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
			var builder strings.Builder
			for j := 0; j < len(p.numsMatrix); j++ {
				char := p.numsMatrix[j][i]
				if char == ' ' {
					continue
				}
				builder.WriteByte(char)
			}
			numAsStr := builder.String()
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
	numberLines, operatorLine := pc.readLines(reader)
	sections := parseOperators(operatorLine)
	pc.createProblems(sections, numberLines)
}

func (pc *ProblemsContainer) readLines(reader io.Reader) ([]string, string) {
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

	return numLines, operatorLine
}

func (pc *ProblemsContainer) createProblems(sections []Section, numberLines []string) {
	for i, section := range sections {
		problem := Problem{
			id:         i,
			numsMatrix: []string{},
			operator:   section.operator,
			start:      section.start,
			end:        section.end,
		}
		for _, numLine := range numberLines {
			problem.numsMatrix = append(problem.numsMatrix, numLine[section.start:section.end])
		}
		pc.problems = append(pc.problems, problem)
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

type Section struct {
	operator string
	start    int
	end      int
}

func parseOperators(operatorLine string) []Section {
	sections := []Section{}
	start := 0
	for i, run := range operatorLine {
		char := string(run)
		isOperator := char == "*" || char == "+"
		isLastChar := i == len(operatorLine)-1
		if isOperator && i > 0 {
			sections = append(sections, Section{
				operator: string(operatorLine[start]),
				start:    start,
				end:      i - 1,
			})
			start = i
		} else if isLastChar {
			sections = append(sections, Section{
				operator: string(operatorLine[start]),
				start:    start,
				end:      i + 1,
			})
		}
	}
	return sections
}
