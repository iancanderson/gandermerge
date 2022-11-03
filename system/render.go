package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/component"
	"github.com/iancanderson/gandermerge/layers"

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
			component.Position,
			component.Sprite,
		)),
}

func (r *render) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	r.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		position := component.GetPosition(entry)
		sprite := component.GetSprite(entry)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.25, 0.25)
		op.GeoM.Translate(position.X, position.Y)
		screen.DrawImage(sprite.Image, op)
	})
}
