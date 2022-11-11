package assets

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/spookypaths/game/assets/images"
	"github.com/iancanderson/spookypaths/game/core"
)

type imageManager struct {
	energyImages map[core.EnergyType]*ebiten.Image
}

// Make sure it conforms to the Manager interface
var _ Manager = (*imageManager)(nil)

var ImageManager = &imageManager{}

func (m *imageManager) Load() {
	m.energyImages = loadEnergyTypeImages()
}

func loadEnergyTypeImages() map[core.EnergyType]*ebiten.Image {
	return map[core.EnergyType]*ebiten.Image{
		core.Electric: LoadImage(images.Electric_png),
		core.Fire:     LoadImage(images.Fire_png),
		core.Ghost:    LoadImage(images.Ghost_png),
		core.Poison:   LoadImage(images.Poison_png),
		core.Psychic:  LoadImage(images.Psychic_png),
	}
}

func EnergyImage(energyType core.EnergyType) *ebiten.Image {
	return ImageManager.energyImages[energyType]
}

func LoadImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
