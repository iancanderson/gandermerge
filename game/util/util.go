package util

import (
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

func BuildGrid(ecs *ecs.ECS) [][]*donburi.Entry {
	query := buildQuery()

	// Store grid in two dimensional array
	grid := make([][]*donburi.Entry, config.Columns)

	// Keep track of where the empty space is
	query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		gridPosition := component.GridPosition.Get(entry)
		if grid[gridPosition.Col] == nil {
			grid[gridPosition.Col] = make([]*donburi.Entry, config.Rows)
		}
		grid[gridPosition.Col][gridPosition.Row] = entry
	})
	return grid
}

func buildQuery() *query.Query {
	return ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.GridPosition,
			component.Sprite,
		))
}

func GridYPosition(row int) float64 {
	return float64(row)*config.RowHeight + config.OrbGridTopMargin
}
