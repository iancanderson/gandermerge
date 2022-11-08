package component

import "github.com/yohamta/donburi"

type HitpointsData struct {
	MaxHitpoints int
	Hitpoints    int
}

var Hitpoints = donburi.NewComponentType[HitpointsData]()

func (h *HitpointsData) IsDead() bool {
	return h.Hitpoints <= 0
}
