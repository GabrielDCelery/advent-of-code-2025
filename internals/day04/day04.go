package day04

import (
	"bufio"
	"context"
	"io"
	"strings"

	"go.uber.org/zap"
)

type RemovalMode int

const (
	RemovalModeSingleLayer RemovalMode = iota
	RemovalModeRecursive
)

type Day4Solver struct {
	logger *zap.Logger
}

func NewDay4Solver(logger *zap.Logger) (*Day4Solver, error) {
	solver := &Day4Solver{
		logger: logger,
	}
	return solver, nil
}

func (d *Day4Solver) Solve(ctx context.Context, reader io.Reader, removalMode RemovalMode) (int, error) {
	grid := &Grid{
		cells: make([][]Cell, 0),
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		cells := strings.Split(line, "")
		row := []Cell{}
		for _, cell := range cells {
			if cell == "@" {
				row = append(row, Roll)
			} else {
				row = append(row, Empty)
			}
		}
		grid.cells = append(grid.cells, row)
	}

	numOfRemovableRools := calculateNumOfRemovableRolls(grid, removalMode, 0)

	return numOfRemovableRools, nil
}

type Cell int

const (
	Empty Cell = iota
	Roll
)

type Grid struct {
	cells [][]Cell
}

func (g *Grid) clone() *Grid {
	cloned := Grid{
		cells: make([][]Cell, 0),
	}
	for _, row := range g.cells {
		clonedRow := []Cell{}
		for _, cell := range row {
			clonedRow = append(clonedRow, cell)
		}
		cloned.cells = append(cloned.cells, clonedRow)
	}
	return &cloned
}

func (g *Grid) removeRolls(coordinates []Coordinate) {
	for _, coordinate := range coordinates {
		g.cells[coordinate.y][coordinate.x] = Empty
	}
}

type Coordinate struct {
	x int
	y int
}

func getRollsReachableViaForklift(grid *Grid) []Coordinate {
	coordinates := make([]Coordinate, 0)
	for cellY, row := range grid.cells {
		for cellX := range row {
			if grid.cells[cellY][cellX] != Roll {
				continue
			}
			neighbouring := 0
			for y := -1; y <= 1; y++ {
				for x := -1; x <= 1; x++ {
					isTheCellItself := x == 0 && y == 0
					if isTheCellItself {
						continue
					}
					neighbourX := cellX + x
					neighbourY := cellY + y
					isInsideGrid := neighbourX >= 0 && neighbourX < len(grid.cells[cellY]) && neighbourY >= 0 && neighbourY < len(grid.cells)
					if !isInsideGrid {
						continue
					}
					neighBourCell := grid.cells[neighbourY][neighbourX]
					if neighBourCell == Roll {
						neighbouring += 1
					}
				}
			}
			if neighbouring < 4 {
				coordinates = append(coordinates, Coordinate{x: cellX, y: cellY})
			}
		}
	}
	return coordinates
}

func calculateNumOfRemovableRolls(grid *Grid, removalMode RemovalMode, sum int) int {
	rollsReachableViaForklift := getRollsReachableViaForklift(grid)
	numOfRollsReachableViaForklift := len(rollsReachableViaForklift)
	if numOfRollsReachableViaForklift == 0 {
		return sum
	}
	sum = sum + numOfRollsReachableViaForklift
	if removalMode == RemovalModeSingleLayer {
		return sum
	}
	cloned := grid.clone()
	cloned.removeRolls(rollsReachableViaForklift)
	return calculateNumOfRemovableRolls(cloned, removalMode, sum)
}
