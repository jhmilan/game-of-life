// Package entities holds the data types for the game of life algorithm
//
//go:generate stringer -type=LifeStatus
//go:generate stringer -type=Transition
package entities

// LifeStatus represents the status of a Cell
type LifeStatus int

const (
	// Alive means a cell is alive
	Alive LifeStatus = iota
	// Dead is the status of a cell when it dies
	Dead
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
func (c Cell) CheckLifeStatus(w World) (LifeStatus, Transition) {
	neighbours := c.GetNeighBours(w)
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
func (c Cell) GetNeighBours(w World) []*Cell {
	cells := []*Cell{}
	for i := c.x - 1; i <= c.x+1; i++ {
		for j := c.y - 1; j <= c.y+1; j++ {
			if w.ValidPosition(i, j) {
				if i == c.x && j == c.y {
					// skip the cell
					continue
				}
				cells = append(cells, w.cells[i][j])
			}
		}
	}

	return cells
}
