package day04

import (
	"bufio"
	"context"
	"io"
	"strings"

	"go.uber.org/zap"
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

func (d *Day4Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
	grid := Grid{
		cells: make([][]Cell, 0),
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		cells := strings.Split(line, "")
		row := []Cell{}
		for _, cell := range cells {
			if cell == "@" {
				row = append(row, Cell(true))
			} else {
				row = append(row, Cell(false))
			}
		}
		grid.cells = append(grid.cells, row)
	}

	sum := 0

	for cellY, row := range grid.cells {
		for cellX := range row {
			if grid.cells[cellY][cellX] != true {
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
					if neighBourCell == true {
						neighbouring += 1
					}
				}
			}
			if neighbouring < 4 {
				sum += 1
			}

		}
	}
	return sum, nil
}

type Cell bool

type Grid struct {
	cells [][]Cell
}
