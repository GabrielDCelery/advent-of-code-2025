package day07

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

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

	beamTracer := NewBeamTracer(teleporter)

	d.logger.Debug("initialised beam tracer", zap.String("beamTracer", fmt.Sprintf("%+v", beamTracer)))

	beamTracer.traceBeamsFromSource()

	splittersCrossedCount := beamTracer.countNumOfSplittersCrossed()
	uniqueBeamsCount := beamTracer.getNumOfBeams()
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

type Beam []Coordinate

func (b *Beam) clone() Beam {
	clone := make(Beam, len(*b))
	copy(clone, *b)
	return clone
}

type BeamTracer struct {
	beams      []Beam
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
		beams:      []Beam{},
		cellMatrix: cellMatrix,
		teleporter: teleporter,
	}
}

func (b *BeamTracer) traceSingleBeam(beam Beam, validBeamsChan chan<- Beam, sem chan struct{}, wg *sync.WaitGroup) {
	sem <- struct{}{}
	defer func() { <-sem }()
	defer wg.Done()

	for {
		head := beam[len(beam)-1]

		hasReachedEnd := head.y == len(b.teleporter.cellMatrix)-1

		if hasReachedEnd {
			validBeamsChan <- beam
			return
		}

		if b.teleporter.cellMatrix[head.y][head.x] == Empty {
			next := Coordinate{x: head.x, y: head.y + 1}
			beam = append(beam, next)
			continue
		}

		if b.teleporter.cellMatrix[head.y][head.x] == Splitter {
			left := Coordinate{x: head.x - 1, y: head.y}
			if left.isWithinBoundaries(&b.teleporter.cellMatrix) {
				wg.Add(1)
				go b.traceSingleBeam(append(beam.clone(), left), validBeamsChan, sem, wg)
			}
			right := Coordinate{x: head.x + 1, y: head.y}
			if right.isWithinBoundaries(&b.teleporter.cellMatrix) {
				wg.Add(1)
				go b.traceSingleBeam(append(beam.clone(), right), validBeamsChan, sem, wg)
			}
			return
		}
	}
}

func (b *BeamTracer) traceBeamsFromSource() {
	var wg sync.WaitGroup
	validBeamsChan := make(chan Beam)
	sem := make(chan struct{}, 5)

	beam := Beam{Coordinate{x: b.teleporter.beamSource.x, y: b.teleporter.beamSource.y}}

	wg.Add(1)
	go b.traceSingleBeam(beam, validBeamsChan, sem, &wg)

	go func() {
		wg.Wait()
		close(validBeamsChan)
	}()

	validBeams := []Beam{}
	for validBeam := range validBeamsChan {
		validBeams = append(validBeams, validBeam)
	}
	b.beams = validBeams
}

func (b *BeamTracer) countNumOfSplittersCrossed() int {
	visitedSplitters := make(map[string]bool, 0)
	for _, beam := range b.beams {
		for _, coordinate := range beam {
			if b.teleporter.cellMatrix[coordinate.y][coordinate.x] == Splitter {
				key := fmt.Sprintf("%d-%d", coordinate.x, coordinate.y)
				visitedSplitters[key] = true
			}
		}
	}
	return len(visitedSplitters)
}

func (b *BeamTracer) getNumOfBeams() int {
	return len(b.beams)
}
