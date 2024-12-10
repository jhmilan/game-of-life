// Package entities holds the data types for the game of life algorithm
//
//go:generate stringer -type=Status
package entities

import (
	"crypto/sha256"
	"fmt"
)

// World represents the board were the cells live
type World struct {
	x     int
	y     int
	cells [][]*Cell
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
func (w World) GetStatus(generations map[string]int) Status {
	numAlive := w.CountAlive()
	if w.CountAlive() == 0 {
		return Extinction
	}
	hash := w.StatusHash()
	if _, ok := generations[hash]; ok {
		return Static
	}
	generations[hash] = numAlive
	return Caos
}

// SetCellAlive sets a cell status as live
func (w World) SetCellAlive(x, y int) {
	if !w.ValidPosition(x, y) {
		panic("invalid cell position")
	}
	w.cells[x][y].status = Alive
}

// Evolve transitions all the cells of the world in one go. It performs a 'step' in time.
// A new instance of the world is returned.
func (w World) Evolve() World {
	newWorld := NewWorld(w.x, w.y)
	for i := range w.cells {
		for j := range w.cells[i] {
			newLifeStatus, _ := w.cells[i][j].CheckLifeStatus(w)
			// we could log the transaction
			if newLifeStatus == Alive {
				newWorld.SetCellAlive(i, j)
			}
		}
	}
	return newWorld
}

// ValidPosition checks if a given x,y position is valid in the current world (not out of bounds)
func (w World) ValidPosition(x, y int) bool {
	if x < 0 || x >= w.x || y < 0 || y >= w.y {
		return false
	}
	return true
}

// StatusHash returns a hash that represents the board status with a fixed length
func (w World) StatusHash() string {
	bitmap := make([]byte, w.x*w.y)
	for i := range w.cells {
		for j := range w.cells[i] {
			if w.cells[i][j].status == Alive {
				bitmap[i*w.x+j] = '1'
			} else {
				bitmap[i*w.x+j] = '0'
			}
		}
	}
	generationID := fmt.Sprintf("%d_%d_%d_%x", w.x, w.y, w.CountAlive(), bitmap)
	sum := sha256.Sum256([]byte(generationID))
	return fmt.Sprintf("%x", sum)
}

// CountAlive return the amount of alive cells
func (w World) CountAlive() int {
	numAlive := 0
	for i := range w.cells {
		for j := range w.cells[i] {
			if w.cells[i][j].status == Alive {
				numAlive++
			}
		}
	}
	return numAlive
}

// Print plots the world in standard output
func (w World) Print() {
	for i := range w.cells {
		for j := range w.cells[i] {
			if w.cells[i][j].status == Dead {
				fmt.Print(" ")
			} else {
				fmt.Print("x")
			}
		}
		fmt.Println()
	}
}

// NewWorld returns a new World instance of the given dimensions. All cells are dead by default
func NewWorld(x, y int) World {
	world := World{
		x: x,
		y: y,
	}
	world.cells = make([][]*Cell, x)
	for i := range world.cells {
		world.cells[i] = make([]*Cell, y)
	}
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			world.cells[i][j] = &Cell{x: i, y: j, status: Dead}
		}
	}
	return world
}
