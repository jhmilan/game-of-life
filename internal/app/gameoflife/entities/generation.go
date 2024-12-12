// Package entities holds the data types for the game of life algorithm
//

package entities

// Generation represents the status of the world at a given moment of time
type Generation struct {
	step     int
	cellsMap CellsMap
}

func newGenerationFromPoints(step int, points [][]int) Generation {
	cellsMap := make(CellsMap, len(points))
	for _, point := range points {
		if len(point) != 2 {
			panic("invalid point - it must contain 2 integers")
		}
		cellsMap[CellPositionFrom(point[0], point[1])] = Cell{
			x:      point[0],
			y:      point[1],
			status: Alive,
		}
	}
	return Generation{
		step:     step,
		cellsMap: cellsMap,
	}
}
