// Package entities holds the data types for the game of life algorithm
//
//go:generate stringer -type=LifeStatus
//go:generate stringer -type=Transition
package entities

import (
	"fmt"
	"maps"
	"slices"
)

// LifeStatus represents the status of a Cell
type LifeStatus int

const (
	// Dead is the status of a cell when it dies
	Dead LifeStatus = iota
	// Alive means a cell is alive
	Alive
)

// Transition represents a Cell transition from one step in time to the next one
type Transition int

const (
	// UnderPopulation happens for any live cell with fewer than two live neighbours. It dies by UnderPopulation
	UnderPopulation Transition = iota
	// LivesOn is the transition for any live cell with two or three live neighbours. It lives on to the next generation
	LivesOn
	// OverPopulation is the destiny for any live cell with more than three live neighbours. It dies by OverPopulation
	OverPopulation
	// Reproduction is what happen when a dead cell is surrounded by exactly three live neighbours. It becomes a live cell
	Reproduction
	// RemainsDead is what happens when a cell is dead and does not rebirth in the next generation
	RemainsDead
)

// CellPosition is a string like 9_22 representing position (9,22) in the board
type CellPosition string

// CellPositionFrom returns a CellPosition given a point (x,y)
func CellPositionFrom(x, y int) CellPosition {
	return CellPosition(fmt.Sprintf("%d_%d", x, y))
}

// CellsMap maps positions to cells (dead or alive)
type CellsMap map[CellPosition]Cell

// GetAliveCells returns a slice with the alive cells
func (cellsMap CellsMap) GetAliveCells() []Cell {
	cells := []Cell{}
	for _, cell := range cellsMap {
		if cell.status == Alive {
			cells = append(cells, cell)
		}
	}
	return cells
}

// GetRelevantCells returns a slice with the cells that must be analyzed
func (cellsMap CellsMap) GetRelevantCells() []Cell {
	aliveCells := cellsMap.GetAliveCells()
	allCells := append([]Cell{}, aliveCells...)
	for _, aliveCell := range aliveCells {
		allCells = append(allCells, aliveCell.GetNeighBours(cellsMap)...)
	}
	newCellMap := make(CellsMap, 0)
	for _, c := range allCells {
		newCellMap[CellPositionFrom(c.x, c.y)] = c
	}
	return slices.Collect(maps.Values(newCellMap))
}

// CellAt returns the cell at position (x,y) if any
func (cellsMap CellsMap) CellAt(x, y int) *Cell {
	if cell, ok := cellsMap[CellPositionFrom(x, y)]; ok {
		return &cell
	}
	return nil
}

// Cell represents a point in the world that can be dead or alive
type Cell struct {
	x      int
	y      int
	status LifeStatus
}

/*
CheckLifeStatus determines the transition from one step in time to the next for the given cell. It depends on the
neighborhood of the cell.
*/
func (c Cell) CheckLifeStatus(cellsMap CellsMap) (LifeStatus, Transition) {
	neighbours := c.GetNeighBours(cellsMap)
	numAliveNeighbors := 0
	for _, neighbor := range neighbours {
		if neighbor.status == Alive {
			numAliveNeighbors++
		}
	}
	if c.status == Alive {
		if numAliveNeighbors < 2 {
			return Dead, UnderPopulation
		}

		if numAliveNeighbors > 3 {
			return Dead, OverPopulation
		}
		return Alive, LivesOn
	}
	if numAliveNeighbors == 3 {
		return Alive, Reproduction
	}
	return Dead, RemainsDead
}

// GetNeighBours return the cells in the neighbourhood
func (c Cell) GetNeighBours(cellsMap CellsMap) []Cell {
	cells := []Cell{}
	for i := c.x - 1; i <= c.x+1; i++ {
		for j := c.y - 1; j <= c.y+1; j++ {
			if i == c.x && j == c.y {
				// skip the cell
				continue
			}
			newCell := Cell{x: i, y: j, status: Dead}
			if cell, ok := cellsMap[CellPositionFrom(i, j)]; ok {
				newCell.status = cell.status
			}
			cells = append(cells, newCell)
		}
	}

	return cells
}
