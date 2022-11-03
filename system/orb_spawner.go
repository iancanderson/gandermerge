package system

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/assets/images"
	"github.com/iancanderson/gandermerge/component"
	"github.com/iancanderson/gandermerge/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const columns = 8
const rows = 8
const columnWidth = 48
const rowHeight = 48

type OrbSpawner struct {
}

func NewOrbSpawner() *OrbSpawner {
	return &OrbSpawner{}
}

func (s *OrbSpawner) Update(ecs *ecs.ECS) {
}

func (s *OrbSpawner) Startup(ecs *ecs.ECS) {
	orbs := ecs.CreateMany(
		layers.LayerOrbs,
		rows*columns,
		component.Position,
		component.Sprite,
	)

	image := loadSprite()

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			entry := ecs.World.Entry(orbs[row*columns+col])

			donburi.SetValue(entry, component.Position,
				component.PositionData{
					X: float64(col) * columnWidth,
					Y: float64(row) * rowHeight,
				})

			donburi.SetValue(entry, component.Sprite,
				component.SpriteData{Image: image})
		}
	}
}

func loadSprite() *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(images.Electric_png))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
