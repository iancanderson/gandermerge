package config

const Columns = 8
const Rows = 8
const ColumnWidth = 96
const RowHeight = 96
const WindowHeight = 1200
const WindowWidth = 768
const OrbGridTopMargin = WindowHeight - (Rows * RowHeight)
const MovesAllowed = 20
const EnemyHitpoints = 140
const AudioSampleRate = 44100
const ModalText = "Welcome to Spooky Paths!\n\n" +
	"Goal: defeat the boss before\nyou run out of moves.\n\n" +
	"Attack: tap and drag to create\na path with the same type of energy,\nrelease to attack.\n\n" +
	"Multipliers: different bosses are\naffected more or less by different\nenergy types:\n\n" +
	"• 2x: boss takes double damage.\n" +
	"• 1x: boss takes normal damage.\n" +
	"• ½x: boss takes half damage.\n\n" +
	"Example: Poison boss (green skull)\ntakes 2x damage from crystal ball orbs,\n½x damage from poison orbs,\nand 1x damage from other orbs.\n\n" +
	"If you attack with 3 crystal ball orbs,\nthe boss will take 6 damage.\n\n" +
	"If you attack with 3 poison orbs,\nthe boss will take 1 damage.\n\n" +
	"If you attack with 3 fire orbs,\nthe boss will take 3 damage.\n\n"
