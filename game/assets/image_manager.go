package assets

import (
	"bytes"
	"image"
	"log"
	"time"

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
	isChristmas := time.Now().Month() == time.December && time.Now().Day() == 25
	isLastThursdayInNovember := time.Now().Month() == time.November &&
		time.Now().Weekday() == time.Thursday &&
		time.Now().Day() >= 24

	if isChristmas {
		m.energyImages = loadChristmasImages()
		return
	}

	if isLastThursdayInNovember {
		m.energyImages = loadThanksgivingImages()
		return
	}

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

func loadThanksgivingImages() map[core.EnergyType]*ebiten.Image {
	return map[core.EnergyType]*ebiten.Image{
		core.Electric: LoadImage(images.Thanksgiving_electric_png),
		core.Fire:     LoadImage(images.Thanksgiving_fire_png),
		core.Ghost:    LoadImage(images.Thanksgiving_ghost_png),
		core.Poison:   LoadImage(images.Thanksgiving_poison_png),
		core.Psychic:  LoadImage(images.Thanksgiving_psychic_png),
	}
}

func loadChristmasImages() map[core.EnergyType]*ebiten.Image {
	return map[core.EnergyType]*ebiten.Image{
		core.Electric: LoadImage(images.Christmas_electric_png),
		core.Fire:     LoadImage(images.Christmas_fire_png),
		core.Ghost:    LoadImage(images.Christmas_ghost_png),
		core.Poison:   LoadImage(images.Christmas_poison_png),
		core.Psychic:  LoadImage(images.Christmas_psychic_png),
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
