package component

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type spriteData struct {
	Image     *ebiten.Image
	X, Y      float64
	scale     float64
	greenTint bool
	redTint   bool
}

var Sprite = donburi.NewComponentType[spriteData]()

func NewSpriteData(image *ebiten.Image, x, y float64) *spriteData {
	return &spriteData{
		Image: image,
		X:     x,
		Y:     y,
		scale: 1.0,
	}
}

// Is there color in the image at the given point?
func (s *spriteData) InColor(x, y int) bool {
	imagePt := s.worldToImage(x, y)
	collideColor := s.Image.At(imagePt.X, imagePt.Y)
	collides := collideColor.(color.RGBA).A > 0

	return collides
}

// Is the given point within the sprite?
func (s *spriteData) In(x, y int) bool {
	imagePt := s.worldToImage(x, y)
	return imagePt.In(s.Image.Bounds())
}

func (s *spriteData) DrawOptions() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s.scale, s.scale)
	op.GeoM.Translate(s.X, s.Y)
	op.Filter = ebiten.FilterLinear
	if s.greenTint {
		op.ColorM.Scale(0.5, 1.0, 0.5, 1)
	} else if s.redTint {
		op.ColorM.Scale(1.0, 0.5, 0.5, 1)
	}
	return op
}

func (s *spriteData) worldToImage(x, y int) image.Point {
	imageX := (float64(x) - s.X) / s.scale
	imageY := (float64(y) - s.Y) / s.scale
	return image.Point{X: int(imageX), Y: int(imageY)}
}

func (s *spriteData) WithScale(scale float64) *spriteData {
	s.scale = scale
	return s
}

func (s *spriteData) WithGreenTint(greenTint bool) *spriteData {
	s.greenTint = greenTint
	return s
}

func (s *spriteData) WithRedTint(redTint bool) *spriteData {
	s.redTint = redTint
	return s
}
