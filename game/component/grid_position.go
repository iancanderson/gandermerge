package component

import (
	"github.com/yohamta/donburi"
)

type GridPositionData struct {
	Row, Col int
}

var GridPosition = donburi.NewComponentType[GridPositionData]()

func (g *GridPositionData) IsAdjacent(other *GridPositionData) bool {
	rowDiff := g.Row - other.Row
	if rowDiff < 0 {
		rowDiff = -rowDiff
	}
	colDiff := g.Col - other.Col
	if colDiff < 0 {
		colDiff = -colDiff
	}
	return rowDiff <= 1 && colDiff <= 1
}
