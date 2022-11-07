package sounds

import (
	_ "embed"
)

var (
	//go:embed ghost-chain.wav
	GhostChain []byte

	//go:embed ghost-merge.wav
	GhostMerge []byte

	//go:embed poison-chain.wav
	PoisonChain []byte

	//go:embed poison-merge.wav
	PoisonMerge []byte
)
