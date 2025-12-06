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

func createProblems(sections []Section, numberLines []string) []Problem {
	problems := []Problem{}
	for _, section := range sections {
		problem := Problem{
			numberRows: []string{},
			operator:   section.operator,
		}
		for _, numLine := range numberLines {
			problem.numberRows = append(problem.numberRows, numLine[section.start:section.end])
		}
		problems = append(problems, problem)
	}
	return problems
}

type Problem struct {
	numberRows []string
	operator   string
}

func (p *Problem) getWidth() int {
	width := 0
	for _, numberRow := range p.numberRows {
		if len(numberRow) > width {
			width = len(numberRow)
		}
	}
	return width
}

func (p *Problem) parseNumbersHorizontally() ([]int, error) {
	nums := []int{}
	for _, numAsStr := range p.numberRows {
		numAsStr = strings.TrimSpace(numAsStr)
		num, err := strconv.Atoi(numAsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid integer '%s'", numAsStr)
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func (p *Problem) parseNumbersVertically() ([]int, error) {
	nums := []int{}
	for i := 0; i < p.getWidth(); i++ {
		var builder strings.Builder
		for j := 0; j < len(p.numberRows); j++ {
			char := p.numberRows[j][i]
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

func (p *Problem) parseNumberRowsToNums(puzzleInterpreter PuzzleInterpreter) ([]int, error) {
	switch puzzleInterpreter {
	case HumanMath:
		return p.parseNumbersHorizontally()
	case CephalopodMath:
		return p.parseNumbersVertically()
	}
	return nil, fmt.Errorf("invalid interpreter '%d'", puzzleInterpreter)
}

func (p *Problem) solve(puzzleInterpreter PuzzleInterpreter) (int, error) {
	numbers, err := p.parseNumberRowsToNums(puzzleInterpreter)
	if err != nil {
		return 0, err
	}
	switch p.operator {
	case "*":
		return multiplyAll(numbers), nil
	case "+":
		return sumAll(numbers), nil
	default:
		return 0, fmt.Errorf("invalid operator %s", p.operator)
	}
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

func multiplyAll(numbers []int) int {
	sum := 1
	for _, number := range numbers {
		sum *= number
	}
	return sum
}

func sumAll(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
