package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/layers"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type render struct {
	query *query.Query
}

var Render = &render{
	query: ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.Sprite,
			component.Selectable,
		)),
}

func (r *render) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	r.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.GetSprite(entry)
		selectable := component.GetSelectable(entry)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(sprite.Scale, sprite.Scale)
		op.GeoM.Translate(sprite.X+4, sprite.Y+4)
		if selectable.Selected {
			op.ColorM.Scale(0.5, 0.5, 0.5, 1)
		}

		screen.DrawImage(sprite.Image, op)
	})
}
