package system

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/spookypaths/game/assets/images"
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/core"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/iancanderson/spookypaths/game/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const spawnYOffset = -150

type OrbSpawner struct {
	images map[core.EnergyType]*ebiten.Image
}

func NewOrbSpawner() *OrbSpawner {
	return &OrbSpawner{}
}

func (s *OrbSpawner) Startup(ecs *ecs.ECS) {
	s.images = loadEnergyTypeImages()

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
	donburi.SetValue(entry, component.Energy,
		component.EnergyData{
			EnergyType: energyType,
		})

	donburi.SetValue(entry, component.Sprite,
		component.SpriteData{
			Image: s.images[energyType],
			X:     4 + float64(col)*config.ColumnWidth,
			Y:     4 + util.GridYPosition(row) + spawnYOffset,
		}.WithScale(0.14).WithGreenTint(energyType == core.Poison).WithRedTint(energyType == core.Fire))

	donburi.SetValue(entry, component.GridPosition,
		component.GridPositionData{
			Row: row,
			Col: col,
		})
}

func loadEnergyTypeImages() map[core.EnergyType]*ebiten.Image {
	return map[core.EnergyType]*ebiten.Image{
		core.Electric: loadImage(images.Electric_png),
		core.Fire:     loadImage(images.Fire_png),
		core.Ghost:    loadImage(images.Ghost_png),
		core.Poison:   loadImage(images.Poison_png),
		core.Psychic:  loadImage(images.Psychic_png),
	}
}

func loadImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
