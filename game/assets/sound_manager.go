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
	bgMusic     *audio.Player
	playSounds  bool
}

// Make sure it conforms to the Manager interface
var _ Manager = (*soundManager)(nil)

var SoundManager = &soundManager{
	chainSounds: make(map[core.EnergyType]*audio.Player),
	mergeSounds: make(map[core.EnergyType]*audio.Player),
	playSounds:  true,
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

	stream, err := vorbis.DecodeWithoutResampling(bytes.NewReader(sounds.BgMusicWithIntro))
	if err != nil {
		panic(err)
	}
	introLength := stream.Length() / 23
	bgStream := audio.NewInfiniteLoopWithIntro(stream, introLength, stream.Length()-introLength)
	s.bgMusic, err = audioContext.NewPlayer(bgStream)
	if err != nil {
		panic(err)
	}
	s.bgMusic.Play()
}

func (s *soundManager) Toggle() {
	if s.bgMusic.IsPlaying() {
		s.bgMusic.Pause()
	} else {
		s.bgMusic.Play()
	}

	s.playSounds = !s.playSounds
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
	if !s.playSounds {
		return
	}

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
	if !s.playSounds {
		return
	}

	player := s.mergeSounds[energyType]
	if player == nil {
		return
	}
	player.Rewind()
	player.Play()
}
