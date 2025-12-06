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
	numberLines, operatorLine := readLines(reader)
	d.logger.Debug("read input to operator line", zap.String("operatorLine", operatorLine))
	d.logger.Debug("read input to number lines", zap.String("numberLines", fmt.Sprintf("%+v", numberLines)))
	sections := parseOperators(operatorLine)
	d.logger.Debug("split operator line to sections", zap.String("sections", fmt.Sprintf("%+v", sections)))
	problems := createProblems(sections, numberLines)
	d.logger.Debug("converted number lines to problems", zap.String("problems", fmt.Sprintf("%+v", problems)))
	solution, err := solveProblems(problems, puzzleInterpreter)
	if err != nil {
		return 0, err
	}
	return solution, nil
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

func readLines(reader io.Reader) ([]string, string) {
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

func createProblems(sections []Section, numberLines []string) []Problem {
	problems := []Problem{}
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
		problems = append(problems, problem)
	}
	return problems
}

func solveProblems(problems []Problem, puzzleInterpreter PuzzleInterpreter) (int, error) {
	sum := 0
	for _, problem := range problems {
		result, err := problem.solve(puzzleInterpreter)
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
