package component

import "github.com/yohamta/donburi"

type ScoreData struct {
	MovesRemaining  int
	TotalEnergyGoal int
}

var Score = donburi.NewComponentType[ScoreData]()

func GetScore(entry *donburi.Entry) *ScoreData {
	return donburi.Get[ScoreData](entry, Score)
}
