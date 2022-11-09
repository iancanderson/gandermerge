package sounds

import (
	_ "embed"
)

var (
	//go:embed fire-chain.ogg
	FireChain []byte
	//go:embed fire-merge.ogg
	FireMerge []byte

	//go:embed ghost-chain.ogg
	GhostChain []byte
	//go:embed ghost-merge.ogg
	GhostMerge []byte

	//go:embed poison-chain.ogg
	PoisonChain []byte
	//go:embed poison-merge.ogg
	PoisonMerge []byte

	//go:embed psychic-chain.ogg
	PsychicChain []byte
	//go:embed psychic-merge.ogg
	PsychicMerge []byte
)
