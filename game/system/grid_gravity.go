package system

import (
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/iancanderson/spookypaths/game/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type gridGravity struct {
	query *query.Query
}

var GridGravity = &gridGravity{
	query: ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.GridPosition,
			component.Sprite,
		)),
}

func (g *gridGravity) Update(ecs *ecs.ECS) {
	grid := util.BuildGrid(ecs)
	g.setNewGridPositions(grid)
	g.animateToGridPositions(ecs)
}

func (g *gridGravity) setNewGridPositions(grid [][]*donburi.Entry) {
	// Move orbs down to fill in gaps
	for col := 0; col < config.Columns; col++ {
		emptyRows := 0
		// Row 0 is the top of the grid
		for row := config.Rows - 1; row >= 0; row-- {
			entry := grid[col][row]
			if entry == nil {
				emptyRows++
			} else {
				// Move the orb down
				gridPosition := component.GridPosition.Get(entry)
				gridPosition.Row += emptyRows
				donburi.Add(entry, component.GridPosition, gridPosition)
			}
		}
	}
}

func (g *gridGravity) animateToGridPositions(ecs *ecs.ECS) {
	// Move orbs towards their grid positions
	g.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.Sprite.Get(entry)
		gridPosition := component.GridPosition.Get(entry)

		if sprite.Y < util.GridYPosition(gridPosition.Row) {
			sprite.Y += 8
		}
	})
}
