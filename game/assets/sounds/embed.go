package sounds

import (
	_ "embed"
)

var (
	//go:embed calvin-ghost-breath.wav
	Ghost_breath_wav []byte

	//go:embed dylan-poison.wav
	Poison_wav []byte

	//go:embed calvin-ian-dylan-psychic-glass.wav
	Psychic_wav []byte
)
