package assets

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type fontManager struct {
	Go36  font.Face
	Go72  font.Face
	Go108 font.Face
}

var FontManager = &fontManager{}

func (f *fontManager) Load() {
	goreg, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	f.Go36 = truetype.NewFace(goreg, &truetype.Options{Size: 36})
	f.Go72 = truetype.NewFace(goreg, &truetype.Options{Size: 72})
	f.Go108 = truetype.NewFace(goreg, &truetype.Options{Size: 108})
}

// Make sure it conforms to the Manager interface
var _ Manager = (*fontManager)(nil)
