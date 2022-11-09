package system

import (
	"math"

	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/layers"
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
		sprite := component.Sprite.Get(entry)
		projectile := component.Projectile.Get(entry)

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

		if math.Abs(sprite.X-projectile.DestX) < projectileVelocity &&
			math.Abs(sprite.Y-projectile.DestY) < projectileVelocity {
			ecs.World.Remove(entry.Entity())
		}
	})
}
