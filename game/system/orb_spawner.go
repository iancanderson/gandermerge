package system

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/game/assets/images"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/iancanderson/gandermerge/game/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const spawnYOffset = 150

type OrbSpawner struct {
	images map[component.EnergyType]*ebiten.Image
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
			X:     float64(col) * config.ColumnWidth,
			Y:     float64(row)*config.RowHeight - spawnYOffset,
			Scale: 0.25,
		})

	donburi.SetValue(entry, component.GridPosition,
		component.GridPositionData{
			Row: row,
			Col: col,
		})
}

func loadEnergyTypeImages() map[component.EnergyType]*ebiten.Image {
	return map[component.EnergyType]*ebiten.Image{
		component.Electric: loadImage(images.Electric_png),
		component.Fire:     loadImage(images.Fire_png),
		component.Ghost:    loadImage(images.Ghost_png),
		component.Poison:   loadImage(images.Poison_png),
		component.Psychic:  loadImage(images.Psychic_png),
	}
}

func loadImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
