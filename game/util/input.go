package util

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputSource interface {
	Position() (int, int)
	JustReleased() bool
}

type MouseInputSource struct{}

func (m MouseInputSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (m MouseInputSource) JustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

type TouchInputSource struct {
	ID ebiten.TouchID
}

func (t TouchInputSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

func (t TouchInputSource) JustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}

func JustPressedInputSource() InputSource {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return MouseInputSource{}
	}

	touchIDs := inpututil.AppendJustPressedTouchIDs([]ebiten.TouchID{})
	if len(touchIDs) > 0 {
		return TouchInputSource{ID: touchIDs[0]}
	}

	return nil
}
