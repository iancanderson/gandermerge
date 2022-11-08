package component

import (
	"image/color"

	"github.com/yohamta/donburi"
	"golang.org/x/image/font"
)

type Bubble int

const (
	BubbleNone = iota
	BubbleLeft
	BubbleRight
	BubbleTop
	BubbleBottom
)

type TextData struct {
	X, Y     int
	FontFace font.Face
	Color    color.Color
	Text     string
	Bubble   Bubble
}

var Text = donburi.NewComponentType[TextData]()
