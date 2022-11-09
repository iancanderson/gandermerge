package component

import (
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/yohamta/donburi"
)

type ScoreData struct {
	MovesRemaining int
	EnemiesAreDead bool
}

func (s *ScoreData) GameOver() bool {
	return s.Won() || s.Lost()
}

func (s *ScoreData) Won() bool {
	return s.EnemiesAreDead
}

func (s *ScoreData) Lost() bool {
	return s.MovesRemaining <= 0 && !s.EnemiesAreDead
}

func (s *ScoreData) NewGame() {
	s.MovesRemaining = config.MovesAllowed
	s.EnemiesAreDead = false
}

var Score = donburi.NewComponentType[ScoreData]()
