package component

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	Image     *ebiten.Image
	X, Y      float64
	Scale     float64
	GreenTint bool
}

var Sprite = donburi.NewComponentType[SpriteData]()

// Is there color in the image at the given point?
func (s *SpriteData) InColor(x, y int) bool {
	imagePt := s.worldToImage(x, y)
	collideColor := s.Image.At(imagePt.X, imagePt.Y)
	collides := collideColor.(color.RGBA).A > 0

	return collides
}

// Is the given point within the sprite?
func (s *SpriteData) In(x, y int) bool {
	imagePt := s.worldToImage(x, y)
	return imagePt.In(s.Image.Bounds())
}

func (s *SpriteData) worldToImage(x, y int) image.Point {
	imageX := (float64(x) - s.X) / s.Scale
	imageY := (float64(y) - s.Y) / s.Scale
	return image.Point{X: int(imageX), Y: int(imageY)}
}

func GetSprite(entry *donburi.Entry) *SpriteData {
	return donburi.Get[SpriteData](entry, Sprite)
}
