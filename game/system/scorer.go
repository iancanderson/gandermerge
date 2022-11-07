package system

import (
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type scorer struct {
}

var Scorer = &scorer{}

func (s *scorer) Startup(ecs *ecs.ECS) {
	entity := ecs.Create(layers.LayerScoreboard, component.Score)
	entry := ecs.World.Entry(entity)

	donburi.SetValue(entry, component.Score, component.ScoreData{
		MovesRemaining: config.MovesAllowed,
		BossHitpoints:  config.EnergyToWin,
	})
}
