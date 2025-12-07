package day07

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"go.uber.org/zap"
)

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

func (d *Day7Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
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
	beamTracer := NewBeamTracer(teleporter)
	d.logger.Debug("initialised beam tracer", zap.String("beamTracer", fmt.Sprintf("%+v", beamTracer)))
	beamTracer.traceFromSource()
	splitterCrossedCount := beamTracer.countNumOfSplittersCrossed()
	d.logger.Debug("counted number of splitters crossed", zap.Int("splitterCrossedCount", splitterCrossedCount))
	return splitterCrossedCount, nil
}

type Coordinate struct {
	x int
	y int
}

func (c Coordinate) isWithinBoundaries(grid *[][]CellType) bool {
	return c.x >= 0 && c.y >= 0 && c.x < len((*grid)[0]) && c.y < len(*grid)
}

func (c Coordinate) getCoordinates() (x int, y int) {
	return c.x, c.y
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

func (t *Teleporter) getSplitterCoordinates() []Coordinate {
	coordinates := []Coordinate{}
	for y, row := range t.cellMatrix {
		for x, cell := range row {
			if cell == Splitter {
				coordinates = append(coordinates, Coordinate{x, y})
			}
		}
	}
	return coordinates
}

type BeamTracer struct {
	cellMatrix [][]bool
	teleporter *Teleporter
}

func NewBeamTracer(teleporter *Teleporter) *BeamTracer {
	height := len(teleporter.cellMatrix)
	width := len((teleporter.cellMatrix)[0])
	cellMatrix := [][]bool{}
	for range height {
		row := make([]bool, width)
		cellMatrix = append(cellMatrix, row)
	}
	return &BeamTracer{
		cellMatrix: cellMatrix,
		teleporter: teleporter,
	}
}

func (b *BeamTracer) _trace(coordinate Coordinate) {
	// the coordinates are out of bounds
	if !coordinate.isWithinBoundaries(&b.teleporter.cellMatrix) {
		return
	}

	x, y := coordinate.getCoordinates()

	// we already visited this cell
	if b.cellMatrix[y][x] == true {
		return
	}

	// mark the cell as visited
	b.cellMatrix[y][x] = true

	// if the visited cell is empty then move down
	if b.teleporter.cellMatrix[y][x] == Empty {
		nextCoordinate := Coordinate{x: x, y: y + 1}
		b._trace(nextCoordinate)
		return
	}

	// if the visited cell is a splitter split the beam
	if b.teleporter.cellMatrix[y][x] == Splitter {
		left := Coordinate{x: x - 1, y: y}
		right := Coordinate{x: x + 1, y: y}
		b._trace(left)
		b._trace(right)
		return
	}
}

func (b *BeamTracer) traceFromSource() {
	b._trace(b.teleporter.beamSource)
}

func (b *BeamTracer) countNumOfSplittersCrossed() int {
	sum := 0
	coordinates := b.teleporter.getSplitterCoordinates()
	for _, coordinate := range coordinates {
		if b.cellMatrix[coordinate.y][coordinate.x] == true {
			sum += 1
		}
	}
	return sum
}
