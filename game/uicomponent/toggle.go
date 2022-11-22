package uicomponent

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/furex/v2"
)

type Toggle struct {
	pressed  bool
	On       bool
	OnToggle func()
	ImageOn  *ebiten.Image
	ImageOff *ebiten.Image
}

var _ furex.ButtonHandler = (*Toggle)(nil)
var _ furex.DrawHandler = (*Toggle)(nil)

func (t *Toggle) HandlePress(x, y int, _ ebiten.TouchID) {
	t.pressed = true
}

func (t *Toggle) HandleRelease(x, y int, isCancel bool) {
	t.pressed = false
	if !isCancel {
		t.On = !t.On
		if t.OnToggle != nil {
			t.OnToggle()
		}
	}
}

func (t *Toggle) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	if t.pressed {
		furex.G.FillRect(screen, &furex.FillRectOpts{
			Rect: frame, Color: color.RGBA{0xaa, 0, 0, 0xff},
		})
	} else {
		furex.G.FillRect(screen, &furex.FillRectOpts{
			Rect: frame, Color: color.RGBA{0, 0xaa, 0, 0xff},
		})
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(frame.Min.X)+4, float64(frame.Min.Y)+4)
	if t.On {
		screen.DrawImage(t.ImageOn, op)
	} else {
		screen.DrawImage(t.ImageOff, op)
	}
}
