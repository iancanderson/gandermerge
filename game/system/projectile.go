package system

import (
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

const projectileVelocity = 20

type projectile struct {
	query *query.Query
}

var Projectile = &projectile{
	query: ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.Projectile,
			component.Sprite,
		),
	),
}

func (p *projectile) Update(ecs *ecs.ECS) {
	p.query.EachEntity(ecs.World, func(entry *donburi.Entry) {

		// Move the projectile
		sprite := component.GetSprite(entry)
		projectile := component.GetProjectile(entry)

		if projectile.DestX > sprite.X {
			sprite.X += projectileVelocity
		} else {
			sprite.X -= projectileVelocity
		}
		if projectile.DestY > sprite.Y {
			sprite.Y += projectileVelocity
		} else {
			sprite.Y -= projectileVelocity
		}
	})
}
