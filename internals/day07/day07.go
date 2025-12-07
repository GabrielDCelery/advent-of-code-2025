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

	explorer := NewBeamExplorer(teleporter)
	splittersCrossedCount, uniqueBeamsCount := explorer.explore()

	d.logger.Debug("ran exploration", zap.String("explorer", fmt.Sprintf("%+v", explorer)))

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

func (c Coordinate) key() string {
	return fmt.Sprintf("%d-%d", c.x, c.y)
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

type SpitterStats struct {
	numOfUniqueBeamsItGenerates int
	hasExploredLeft             bool
	hasExploredRight            bool
}

type BeamExplorer struct {
	teleporter       *Teleporter
	head             Coordinate
	splittersVisited []Coordinate
	splitterStatsMap map[string]SpitterStats
}

func NewBeamExplorer(teleporter *Teleporter) *BeamExplorer {
	return &BeamExplorer{
		teleporter:       teleporter,
		head:             Coordinate{x: teleporter.beamSource.x, y: teleporter.beamSource.y},
		splittersVisited: []Coordinate{},
		splitterStatsMap: make(map[string]SpitterStats, 0),
	}
}

func (b *BeamExplorer) moveHeadDown() {
	b.head = Coordinate{x: b.head.x, y: b.head.y + 1}
}

func (b *BeamExplorer) moveHeadLeft() {
	b.head = Coordinate{x: b.head.x - 1, y: b.head.y}
}

func (b *BeamExplorer) moveHeadRight() {
	b.head = Coordinate{x: b.head.x + 1, y: b.head.y}
}

func (b *BeamExplorer) moveHeadBackToLastFullyUnexploredSplitter() {
	for {
		if len(b.splittersVisited) == 0 {
			b.head = Coordinate{x: b.teleporter.beamSource.x, y: b.teleporter.beamSource.y}
			return
		}
		last := b.splittersVisited[len(b.splittersVisited)-1]
		key := last.key()
		if b.splitterStatsMap[key].hasExploredLeft && b.splitterStatsMap[key].hasExploredRight {
			b.splittersVisited = b.splittersVisited[:len(b.splittersVisited)-1]
			continue
		} else {
			b.head = Coordinate{x: last.x, y: last.y}
			return
		}
	}
}

func (b *BeamExplorer) hasReachedEndOfTeleporter() bool {
	return b.head.y == len(b.teleporter.cellMatrix)-1
}

func (b *BeamExplorer) getHeadCellType() CellType {
	return b.teleporter.cellMatrix[b.head.y][b.head.x]
}

func (b *BeamExplorer) isHeadInsideTeleporter() bool {
	return b.head.isWithinBoundaries(&b.teleporter.cellMatrix)
}

func (b *BeamExplorer) incrementHowManyUniqueBeamsSplittersCanGenerate() {
	for _, splitter := range b.splittersVisited {
		key := splitter.key()
		splitterStats := b.splitterStatsMap[key]
		b.splitterStatsMap[key] = SpitterStats{
			hasExploredLeft:             splitterStats.hasExploredLeft,
			hasExploredRight:            splitterStats.hasExploredRight,
			numOfUniqueBeamsItGenerates: splitterStats.numOfUniqueBeamsItGenerates + 1,
		}
	}
}

func (b *BeamExplorer) explore() (int, int) {
	uniqueBeamsCount := 0

	for {
		// if the current path we are exploring made us leave the teleporter we have to go back
		if !b.isHeadInsideTeleporter() {
			b.moveHeadBackToLastFullyUnexploredSplitter()
			continue
		}

		// it the current beam reached the end of the teleporter we finshed tracing it
		if b.hasReachedEndOfTeleporter() {
			uniqueBeamsCount += 1
			b.incrementHowManyUniqueBeamsSplittersCanGenerate()
			// move back to the head to the last unexplored splitter so we can continue exploring
			b.moveHeadBackToLastFullyUnexploredSplitter()
			// if there are no more unexplored splitters we have finished
			if len(b.splittersVisited) == 0 {
				// Check if we're back at the source and the first splitter is fully explored
				if b.head.x == b.teleporter.beamSource.x && b.head.y == b.teleporter.beamSource.y {
					// Move down to first splitter and check if fully explored
					firstSplitterKey := fmt.Sprintf("%d-%d", 7, 2)
					if stats, ok := b.splitterStatsMap[firstSplitterKey]; ok && stats.hasExploredLeft && stats.hasExploredRight {
						break
					}
				}
			}
			continue
		}

		if b.getHeadCellType() == Empty {
			b.moveHeadDown()
			continue
		}

		if b.getHeadCellType() == Splitter {
			key := b.head.key()
			splitter, ok := b.splitterStatsMap[key]

			if !ok {
				splitter = SpitterStats{
					hasExploredLeft:             false,
					hasExploredRight:            false,
					numOfUniqueBeamsItGenerates: 0,
				}
			}

			if splitter.hasExploredLeft && splitter.hasExploredRight {
				// This splitter was fully explored in a previous beam
				// Add its cached beam count to our total
				uniqueBeamsCount += splitter.numOfUniqueBeamsItGenerates
				// But DON'T increment parent splitters - that causes overcounting
				b.moveHeadBackToLastFullyUnexploredSplitter()
				if len(b.splittersVisited) == 0 {
					break
				}
			} else if splitter.hasExploredLeft {
				splitter.hasExploredRight = true
				b.splitterStatsMap[key] = splitter
				b.splittersVisited = append(b.splittersVisited, Coordinate{x: b.head.x, y: b.head.y})
				b.moveHeadRight()
			} else {
				splitter.hasExploredLeft = true
				b.splitterStatsMap[key] = splitter
				b.splittersVisited = append(b.splittersVisited, Coordinate{x: b.head.x, y: b.head.y})
				b.moveHeadLeft()
			}
		}
	}
	return len(b.splitterStatsMap), uniqueBeamsCount
}
