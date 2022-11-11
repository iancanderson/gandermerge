package system

import (
	"github.com/iancanderson/spookypaths/game/assets"
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/core"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/iancanderson/spookypaths/game/util"
	"github.com/yohamta/donburi/ecs"
)

const spawnYOffset = -150

type OrbSpawner struct{}

func NewOrbSpawner() *OrbSpawner {
	return &OrbSpawner{}
}

func (s *OrbSpawner) Startup(ecs *ecs.ECS) {
	for row := 0; row < config.Rows; row++ {
		for col := 0; col < config.Columns; col++ {
			s.spawnOrb(ecs, col, row)
		}
	}
}

func (s *OrbSpawner) Update(ecs *ecs.ECS) {
	// Spawn orbs to get to the desired number
	grid := util.BuildGrid(ecs)

	for col := 0; col < config.Columns; col++ {
		for row := 0; row < config.Rows; row++ {
			if grid[col][row] == nil {
				s.spawnOrb(ecs, col, row)
			}
		}
	}
}

func (s *OrbSpawner) spawnOrb(ecs *ecs.ECS, col, row int) {
	orb := ecs.Create(
		layers.LayerOrbs,
		component.Energy,
		component.GridPosition,
		component.Selectable,
		component.Sprite,
	)
	entry := ecs.World.Entry(orb)

	energyType := component.RandomEnergyType()
	component.Energy.Set(entry, &component.EnergyData{
		EnergyType: energyType,
	})

	component.Sprite.Set(entry, component.NewSpriteData(
		assets.EnergyImage(energyType),
		4+float64(col)*config.ColumnWidth,
		4+util.GridYPosition(row)+spawnYOffset,
	).WithScale(0.14).WithGreenTint(energyType == core.Poison).WithRedTint(energyType == core.Fire))

	component.GridPosition.Set(entry, &component.GridPositionData{
		Row: row,
		Col: col,
	})
}
