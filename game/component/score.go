package component

import (
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/yohamta/donburi"
)

type ScoreData struct {
	MovesRemaining int
	EnemiesAreDead bool
}

func (s *ScoreData) IsGameOver() bool {
	return s.MovesRemaining <= 0 || s.EnemiesAreDead
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

func GetScore(entry *donburi.Entry) *ScoreData {
	return donburi.Get[ScoreData](entry, Score)
}
