package assets

import (
	"github.com/golang/freetype/truetype"
	"github.com/iancanderson/spookypaths/game/assets/fonts"
	"golang.org/x/image/font"
)

type fontManager struct {
	Creepster72  font.Face
	Creepster160 font.Face
	Creepster200 font.Face
	Mona36       font.Face
	Mona72       font.Face
	Mona108      font.Face
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

	mona, err := truetype.Parse(fonts.Mona_Sans)
	if err != nil {
		panic(err)
	}
	f.Mona36 = truetype.NewFace(mona, &truetype.Options{Size: 36})
	f.Mona72 = truetype.NewFace(mona, &truetype.Options{Size: 72})
	f.Mona108 = truetype.NewFace(mona, &truetype.Options{Size: 72})
}

// Make sure it conforms to the Manager interface
var _ Manager = (*fontManager)(nil)
