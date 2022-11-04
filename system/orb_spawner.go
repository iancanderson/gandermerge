package system

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/assets/images"
	"github.com/iancanderson/gandermerge/component"
	"github.com/iancanderson/gandermerge/config"
	"github.com/iancanderson/gandermerge/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const desiredOrbCount = config.Rows * config.Columns

type OrbSpawner struct {
}

func NewOrbSpawner() *OrbSpawner {
	return &OrbSpawner{}
}

func (s *OrbSpawner) Update(ecs *ecs.ECS) {
	// Spawn orbs to get to the desired number
}

func (s *OrbSpawner) Startup(ecs *ecs.ECS) {
	orbs := ecs.CreateMany(
		layers.LayerOrbs,
		desiredOrbCount,
		component.Energy,
		component.GridPosition,
		component.Selectable,
		component.Sprite,
	)

	images := loadEnergyTypeImages()

	for row := 0; row < config.Rows; row++ {
		for col := 0; col < config.Columns; col++ {
			entry := ecs.World.Entry(orbs[row*config.Columns+col])

			energyType := component.RandomEnergyType()
			donburi.SetValue(entry, component.Energy,
				component.EnergyData{
					EnergyType: energyType,
				})

			donburi.SetValue(entry, component.Sprite,
				component.SpriteData{
					Image: images[energyType],
					X:     float64(col) * config.ColumnWidth,
					Y:     float64(row) * config.RowHeight,
					Scale: 0.25,
				})

			donburi.SetValue(entry, component.GridPosition,
				component.GridPositionData{
					Row: row,
					Col: col,
				})
		}
	}
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
