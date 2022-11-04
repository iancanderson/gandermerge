package system

import (
	"github.com/iancanderson/gandermerge/component"
	"github.com/iancanderson/gandermerge/config"
	"github.com/iancanderson/gandermerge/layers"
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
	grid := g.buildGrid(ecs)
	g.setNewGridPositions(grid)
	g.animateToGridPositions(ecs)
}

func (g *gridGravity) buildGrid(ecs *ecs.ECS) [][]*donburi.Entry {
	// Store grid in two dimensional array
	grid := make([][]*donburi.Entry, config.Columns)

	// Keep track of where the empty space is
	g.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		gridPosition := component.GetGridPosition(entry)
		if grid[gridPosition.Col] == nil {
			grid[gridPosition.Col] = make([]*donburi.Entry, config.Rows)
		}
		grid[gridPosition.Col][gridPosition.Row] = entry
	})
	return grid
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
				gridPosition := component.GetGridPosition(entry)
				gridPosition.Row += emptyRows
				donburi.Add(entry, component.GridPosition, gridPosition)
			}
		}
	}
}

func (g *gridGravity) animateToGridPositions(ecs *ecs.ECS) {
	// Move orbs towards their grid positions
	g.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.GetSprite(entry)
		gridPosition := component.GetGridPosition(entry)

		if sprite.Y < float64(gridPosition.Row)*config.RowHeight {
			sprite.Y += 8
		}
	})
}
