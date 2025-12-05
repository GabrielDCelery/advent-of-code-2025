package day04

import (
	"bufio"
	"context"
	"io"
	"strings"

	"go.uber.org/zap"
)

const MinNeighborsForInaccessible = 4

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
		cells: make([][]Cell, len(g.cells)),
	}
	for y, row := range g.cells {
		cloned.cells[y] = make([]Cell, len(row))
		copy(cloned.cells[y], row)
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

func (c Coordinate) isInsideGrid(grid *Grid) bool {
	return c.x >= 0 && c.x < len(grid.cells[0]) && c.y >= 0 && c.y < len(grid.cells)
}

func (c Coordinate) getNeighbours() []Coordinate {
	neighbours := make([]Coordinate, 0)
	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			isTheCellItself := x == 0 && y == 0
			if isTheCellItself {
				continue
			}
			neighbour := Coordinate{x: c.x + x, y: c.y + y}
			neighbours = append(neighbours, neighbour)
		}
	}
	return neighbours
}

func getRollsReachableViaForklift(grid *Grid) []Coordinate {
	coordinates := make([]Coordinate, 0)
	for cellY, row := range grid.cells {
		for cellX := range row {
			cell := Coordinate{x: cellX, y: cellY}
			if grid.cells[cell.y][cell.x] != Roll {
				continue
			}
			rollsNextToCell := 0
			neighbours := cell.getNeighbours()
			for _, neighbour := range neighbours {
				if !neighbour.isInsideGrid(grid) {
					continue
				}
				if grid.cells[neighbour.y][neighbour.x] != Roll {
					continue
				}
				rollsNextToCell += 1
			}
			if rollsNextToCell < MinNeighborsForInaccessible {
				coordinates = append(coordinates, cell)
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
