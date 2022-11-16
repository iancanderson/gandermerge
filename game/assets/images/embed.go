package images

import (
	_ "embed"
	_ "image/png"
)

var (
	//go:embed electric.png
	Electric_png []byte

	//go:embed fire.png
	Fire_png []byte

	//go:embed ghost.png
	Ghost_png []byte

	//go:embed information.png
	Information_png []byte

	//go:embed thanksgiving-electric.png
	Thanksgiving_electric_png []byte

	//go:embed thanksgiving-fire.png
	Thanksgiving_fire_png []byte

	//go:embed thanksgiving-ghost.png
	Thanksgiving_ghost_png []byte

	//go:embed thanksgiving-poison.png
	Thanksgiving_poison_png []byte

	//go:embed thanksgiving-psychic.png
	Thanksgiving_psychic_png []byte

	//go:embed poison.png
	Poison_png []byte

	//go:embed psychic.png
	Psychic_png []byte
)
