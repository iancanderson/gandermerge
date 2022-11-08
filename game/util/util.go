package util

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/layers"
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

type InputSource interface {
	Position() (int, int)
	JustReleased() bool
}

type MouseInputSource struct{}

func (m MouseInputSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (m MouseInputSource) JustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

type TouchInputSource struct {
	ID ebiten.TouchID
}

func (t TouchInputSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

func (t TouchInputSource) JustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}

func JustPressedInputSource() InputSource {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return MouseInputSource{}
	}

	touchIDs := inpututil.AppendJustPressedTouchIDs([]ebiten.TouchID{})
	if len(touchIDs) > 0 {
		return TouchInputSource{ID: touchIDs[0]}
	}

	return nil
}
