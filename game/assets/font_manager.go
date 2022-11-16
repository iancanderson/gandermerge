package assets

import (
	"github.com/golang/freetype/truetype"
	"github.com/iancanderson/spookypaths/game/assets/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type fontManager struct {
	Creepster72  font.Face
	Creepster160 font.Face
	Creepster200 font.Face
	Go36         font.Face
	Go72         font.Face
	Go108        font.Face
}

var FontManager = &fontManager{}

func (f *fontManager) Load() {
	creepsterreg, err := truetype.Parse(fonts.Creepster_regular)
	if err != nil {
		panic(err)
	}
	f.Creepster72 = truetype.NewFace(creepsterreg, &truetype.Options{Size: 72})
	f.Creepster160 = truetype.NewFace(creepsterreg, &truetype.Options{Size: 160})
	f.Creepster200 = truetype.NewFace(creepsterreg, &truetype.Options{Size: 200})

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
