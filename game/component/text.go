package component

import (
	"image/color"

	"github.com/yohamta/donburi"
	"golang.org/x/image/font"
)

type TextData struct {
	X, Y     int
	FontFace font.Face
	Color    color.Color
	Text     string
}

var Text = donburi.NewComponentType[TextData]()
