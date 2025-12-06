package day06

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
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

func (d *Day6Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
	scanner := bufio.NewScanner(reader)
	problemContainer := NewProblemContainer()

	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`\s+`)
		line = re.ReplaceAllString(strings.TrimSpace(line), " ")
		d.logger.Debug("cleaned up input", zap.String("line", line))
		components := strings.Split(line, " ")
		if !problemContainer.isInitialised() {
			problemContainer.init(len(components))
		}
		isLastLine := strings.ContainsAny(components[0], `+*`)
		if isLastLine {
			for problemID, operator := range components {
				problemContainer.addOpertor(problemID, operator)
			}
		} else {
			for problemID, numberStr := range components {
				number, err := strconv.Atoi(numberStr)
				if err != nil {
					return 0, fmt.Errorf("value '%s' can not be converted to integer", numberStr)
				}
				problemContainer.addNumber(problemID, number)
			}
		}
	}
	d.logger.Debug("parsed problems into container", zap.String("problems", fmt.Sprintf("%+v", problemContainer)))
	sum := 0
	for _, problem := range problemContainer.problems {
		if problem.operator == "*" {
			miniSum := 1
			for _, number := range problem.numbers {
				miniSum = miniSum * number
			}
			sum += miniSum
		}
		if problem.operator == "+" {
			miniSum := 0
			for _, number := range problem.numbers {
				miniSum = miniSum + number
			}
			sum += miniSum
		}
	}
	return sum, nil
}

type Problem struct {
	id       int
	numbers  []int
	operator string
}

type ProblemContainer struct {
	problems []Problem
}

func NewProblemContainer() *ProblemContainer {
	return &ProblemContainer{
		problems: []Problem{},
	}
}

func (p *ProblemContainer) isInitialised() bool {
	return len(p.problems) != 0
}

func (p *ProblemContainer) init(numOfProblems int) {
	p.problems = make([]Problem, 0, numOfProblems)
	for problemID := range numOfProblems {
		p.problems = append(p.problems, Problem{
			id:       problemID,
			numbers:  []int{},
			operator: ""},
		)
	}
}

func (p *ProblemContainer) addNumber(problemID int, number int) {
	p.problems[problemID].numbers = append(p.problems[problemID].numbers, number)
}

func (p *ProblemContainer) addOpertor(problemID int, operator string) {
	p.problems[problemID].operator = operator
}
