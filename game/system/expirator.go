package system

import (
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type expirator struct {
	query *query.Query
}

var Expirator = &expirator{
	query: ecs.NewQuery(
		layers.LayerEnemy,
		filter.Contains(
			component.Expiration,
		)),
}

func (e *expirator) Update(ecs *ecs.ECS) {
	e.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		expiration := component.Expiration.Get(entry)
		expiration.TTL -= ecs.Time.DeltaTime()
		if expiration.TTL <= 0 {
			ecs.World.Remove(entry.Entity())
		}
	})
}
