package component

import (
	"math/rand"

	"github.com/yohamta/donburi"
)

type EnergyType int

const (
	Electric EnergyType = iota
	Fire
	Ghost
	Poison
	Psychic
)

func RandomEnergyType() EnergyType {
	return EnergyType(rand.Intn(int(Psychic) + 1))
}

type EnergyData struct {
	EnergyType EnergyType
}

var Energy = donburi.NewComponentType[EnergyData]()
