package uicomponent

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/spookypaths/game/assets"
	"github.com/yohamta/furex/v2"
)

type Button struct {
	pressed bool
	Text    string
	OnClick func()
}

var _ furex.ButtonHandler = (*Button)(nil)
var _ furex.DrawHandler = (*Button)(nil)

func (b *Button) HandlePress(x, y int, t ebiten.TouchID) {
	b.pressed = true
}

func (b *Button) HandleRelease(x, y int, isCancel bool) {
	b.pressed = false
	if !isCancel {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}

func (b *Button) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	if b.pressed {
		furex.G.FillRect(screen, &furex.FillRectOpts{
			Rect: frame, Color: color.RGBA{0xaa, 0, 0, 0xff},
		})
	} else {
		furex.G.FillRect(screen, &furex.FillRectOpts{
			Rect: frame, Color: color.RGBA{0, 0xaa, 0, 0xff},
		})
	}

	fontface := assets.FontManager.Creepster72
	textBounds := text.BoundString(fontface, b.Text)
	xOffset := frame.Dx()/2 - textBounds.Dx()/2
	yOffset := frame.Dy()/2 + textBounds.Dy()/2

	text.Draw(screen, b.Text, fontface, frame.Min.X+xOffset, frame.Min.Y+yOffset, color.White)
}
