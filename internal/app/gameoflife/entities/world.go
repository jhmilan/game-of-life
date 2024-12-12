// Package entities holds the data types for the game of life algorithm
//
//go:generate stringer -type=Status
package entities

import (
	"fmt"
)

// World represents the board were the cells live
type World struct {
	x           int
	y           int
	generations []Generation
}

// NewWorld returns a new World instance of the given dimensions. All cells are dead by default
func NewWorld(x, y int, points [][]int) World {
	world := World{
		x:           x,
		y:           y,
		generations: []Generation{newGenerationFromPoints(0, points)},
	}

	return world
}

// Status states the world status at some point
type Status int

const (
	// Caos is an apparent random evolution. No extinction no defined stabilization
	Caos Status = iota
	// Extinction happens when all the cells are dead
	Extinction
	// Static is a configuration where all the structures remain the same generation after generation
	Static
	// Oscillator is a stabilization where patterns are periodically repeated
	Oscillator
	// InfiniteGrowth population grows getting more and more space in the board  (world)
	InfiniteGrowth
	// LimitedGrowth population grows till a limit, then is moves around (spaceships)
	LimitedGrowth
)

// GetStatus returns the status of the world
func (w World) GetStatus() Status {
	numAlive := w.CountAlive()
	if numAlive == 0 {
		return Extinction
	}
	// hash := w.StatusHash()
	// if _, ok := generations[hash]; ok {
	// 	return Static
	// }
	// generations[hash] = numAlive
	return Caos
}

// LastGeneration returns the last generation of cells
func (w World) LastGeneration() Generation {
	if len(w.generations) == 0 {
		panic("empty world - no generations")
	}
	return w.generations[len(w.generations)-1]
}

// Step evolves the world. It creates a new generation
func (w *World) Step() {
	lastGeneration := w.LastGeneration()
	newCellsMap := make(CellsMap, 0)

	for _, cell := range lastGeneration.cellsMap.GetRelevantCells() {
		status, _ := cell.CheckLifeStatus(lastGeneration.cellsMap)
		cell.status = status
		newCellsMap[CellPositionFrom(cell.x, cell.y)] = cell
	}

	w.generations = append(w.generations, Generation{
		step:     lastGeneration.step + 1,
		cellsMap: newCellsMap,
	})
}

// CountAlive returns the amount of alive cells
func (w World) CountAlive() int {
	return len(w.LastGeneration().cellsMap.GetAliveCells())
}

// GetGenerationStep returns the generation step
func (w World) GetGenerationStep() int {
	return w.LastGeneration().step
}

// Print plots the world in standard output
func (w World) Print() {
	cellsMap := w.LastGeneration().cellsMap
	for x := range w.x {
		for y := range w.y {
			cell := cellsMap.CellAt(x, y)
			if cell == nil || cell.status == Dead {
				fmt.Print(" ")
			} else {
				fmt.Print("x")
			}
		}
		fmt.Println()
	}
}
