package core

type EnergyType int

const (
	Electric EnergyType = iota
	Fire
	Ghost
	Poison
	Psychic
)

func ScaleAttack(energyAmount int, attackType EnergyType, defenseType EnergyType) int {
	return int(float64(energyAmount) * attackMultiplier(attackType, defenseType))
}

// - fire
//   - immune to fire
//   - weak to poison
//
// - electric
//   - immune to electric
//   - weak to fire
//
// - poison
//   - immune to poison
//   - weak to psychic
//
// - psychic
//   - immune to ghosts
//   - weak to electric
//
// - ghost
//   - immune to psychic
//   - weak to ghosts
func attackMultiplier(attackType EnergyType, defenseType EnergyType) float64 {
	switch defenseType {
	case Fire:
		switch attackType {
		case Fire:
			return 0.5
		case Poison:
			return 2
		}
	case Electric:
		switch attackType {
		case Electric:
			return 0.5
		case Fire:
			return 2
		}
	case Poison:
		switch attackType {
		case Poison:
			return 0.5
		case Psychic:
			return 2
		}
	case Psychic:
		switch attackType {
		case Ghost:
			return 0.5
		case Electric:
			return 2
		}
	case Ghost:
		switch attackType {
		case Psychic:
			return 0.5
		case Ghost:
			return 2
		}
	}

	return 1
}
