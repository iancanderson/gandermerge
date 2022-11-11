package assets

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/iancanderson/spookypaths/game/assets/sounds"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/core"
)

type soundManager struct {
	chainSounds map[core.EnergyType]*audio.Player
	mergeSounds map[core.EnergyType]*audio.Player
}

// Make sure it conforms to the Manager interface
var _ Manager = (*soundManager)(nil)

var SoundManager = &soundManager{
	chainSounds: make(map[core.EnergyType]*audio.Player),
	mergeSounds: make(map[core.EnergyType]*audio.Player),
}

func (s *soundManager) Load() {
	audioContext := audio.NewContext(config.AudioSampleRate)

	chainSoundData := map[core.EnergyType][]byte{
		core.Electric: sounds.ElectricChain,
		core.Fire:     sounds.FireChain,
		core.Ghost:    sounds.GhostChain,
		core.Poison:   sounds.PoisonChain,
		core.Psychic:  sounds.PsychicChain,
	}

	s.chainSounds = loadSounds(chainSoundData, audioContext)

	mergeSoundData := map[core.EnergyType][]byte{
		core.Electric: sounds.ElectricMerge,
		core.Fire:     sounds.FireMerge,
		core.Ghost:    sounds.GhostMerge,
		core.Poison:   sounds.PoisonMerge,
		core.Psychic:  sounds.PsychicMerge,
	}
	s.mergeSounds = loadSounds(mergeSoundData, audioContext)
}

func loadSounds(soundData map[core.EnergyType][]byte, audioContext *audio.Context) map[core.EnergyType]*audio.Player {
	sounds := make(map[core.EnergyType]*audio.Player)

	for energyType, data := range soundData {
		stream, err := vorbis.DecodeWithoutResampling(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		player, err := audioContext.NewPlayer(stream)
		if err != nil {
			panic(err)
		}
		sounds[energyType] = player
	}
	return sounds
}

func (s *soundManager) PlayChain(energyType core.EnergyType) {
	player := s.chainSounds[energyType]
	if player == nil {
		return
	}
	player.Rewind()
	player.Play()
}

func (s *soundManager) PauseChain(energyType core.EnergyType) {
	player := s.chainSounds[energyType]
	if player == nil {
		return
	}
	player.Pause()
}

func (s *soundManager) PlayMerge(energyType core.EnergyType) {
	player := s.mergeSounds[energyType]
	if player == nil {
		return
	}
	player.Rewind()
	player.Play()
}
