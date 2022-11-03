package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/iancanderson/gandermerge/component"
	"github.com/iancanderson/gandermerge/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type input struct {
	selectableQuery *query.Query
}

var Input = &input{
	selectableQuery: ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.Position,
			component.Selectable,
		)),
}

//TODO: see example for input system: https://github.com/hajimehoshi/ebiten/blob/main/examples/drag/main.go

type InputSource interface {
	Position() (int, int)
}

type MouseInputSource struct{}

func (m *MouseInputSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (r *input) Update(ecs *ecs.ECS) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_ = MouseInputSource{}
		// query for selectable entities

		r.selectableQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
		})
	}
}
