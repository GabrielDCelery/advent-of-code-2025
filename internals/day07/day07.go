package day07

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"go.uber.org/zap"
)

type Solution struct {
	SplittersCrossedCount int
	UniqueBeamsCount      int
}

type Day7Solver struct {
	logger *zap.Logger
}

func NewDay7Solver(logger *zap.Logger) (*Day7Solver, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	solver := &Day7Solver{
		logger: logger,
	}
	return solver, nil
}

func (d *Day7Solver) Solve(ctx context.Context, reader io.Reader) (Solution, error) {
	teleporter := NewTeleporter()

	scanner := bufio.NewScanner(reader)

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		cellMatrix := strings.Split(line, "")
		row := []CellType{}
		for x, cell := range cellMatrix {
			switch cell {
			case "S":
				teleporter.beamSource = Coordinate{x, y}
				row = append(row, Empty)
			case ".":
				row = append(row, Empty)
			case "^":
				row = append(row, Splitter)
			}
		}
		teleporter.cellMatrix = append(teleporter.cellMatrix, row)
		y++
	}

	d.logger.Debug("generated teleporter diagram", zap.String("teleporter", fmt.Sprintf("%+v", teleporter)))

	splittersCrossedCount, uniqueBeamsCount := traceBeams(teleporter)

	return Solution{
		SplittersCrossedCount: splittersCrossedCount,
		UniqueBeamsCount:      uniqueBeamsCount,
	}, nil
}

type Coordinate struct {
	x int
	y int
}

func (c Coordinate) isWithinBoundaries(grid *[][]CellType) bool {
	return c.x >= 0 && c.y >= 0 && c.x < len((*grid)[0]) && c.y < len(*grid)
}

type CellType int

const (
	Empty CellType = iota
	Splitter
)

type Teleporter struct {
	beamSource Coordinate
	cellMatrix [][]CellType
}

func NewTeleporter() *Teleporter {
	cellMatrix := make([][]CellType, 0)
	return &Teleporter{
		beamSource: Coordinate{},
		cellMatrix: cellMatrix,
	}
}

func traceBeams(teleporter *Teleporter) (int, int) {
	uniqueSplittersVisited := make(map[string]bool, 0)
	beamCount := 0

	stack := []Coordinate{{x: teleporter.beamSource.x, y: teleporter.beamSource.y}}

	for len(stack) > 0 {
		current := stack[len(stack)-1]

		// if the current beam reached the end of the teleporter we finished tracing it
		if current.y == len(teleporter.cellMatrix)-1 {
			beamCount += 1
			stack = stack[:len(stack)-1]
			continue
		}

		// if the current beam is on an empty cell move it down
		if teleporter.cellMatrix[current.y][current.x] == Empty {
			next := Coordinate{x: current.x, y: current.y + 1}
			stack = append(stack[:len(stack)-1], next)
			continue
		}

		// if we are on a splitter mark the splitter visited and split the beam
		if teleporter.cellMatrix[current.y][current.x] == Splitter {
			key := fmt.Sprintf("%d-%d", current.x, current.y)
			uniqueSplittersVisited[key] = true
			left := Coordinate{x: current.x - 1, y: current.y}
			right := Coordinate{x: current.x + 1, y: current.y}
			stack = stack[:len(stack)-1]
			if left.isWithinBoundaries(&teleporter.cellMatrix) {
				stack = append(stack, left)
			}
			if right.isWithinBoundaries(&teleporter.cellMatrix) {
				stack = append(stack, right)
			}
		}
	}
	return len(uniqueSplittersVisited), beamCount
}
