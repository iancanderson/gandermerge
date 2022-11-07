package component

import (
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/yohamta/donburi"
)

type ScoreData struct {
	MovesRemaining int
	BossHitpoints  int
}

func (s *ScoreData) IsGameOver() bool {
	return s.MovesRemaining <= 0 || s.BossHitpoints <= 0
}

func (s *ScoreData) Won() bool {
	return s.BossHitpoints <= 0
}

func (s *ScoreData) Lost() bool {
	return s.MovesRemaining <= 0
}

func (s *ScoreData) NewGame() {
	s.MovesRemaining = config.MovesAllowed
	s.BossHitpoints = config.EnergyToWin
}

var Score = donburi.NewComponentType[ScoreData]()

func GetScore(entry *donburi.Entry) *ScoreData {
	return donburi.Get[ScoreData](entry, Score)
}