package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type game struct {
	screen ebiten.Game
}

func NewGame() *game {
	game := game{}
	screen := NewMainMenuScreen(
		game.startDailyLevel,
		game.startRandomLevel,
	)
	game.screen = screen
	return &game
}

func (g *game) Update() error {
	err := g.screen.Update()
	return err
}

func (g *game) Draw(screen *ebiten.Image) {
	g.screen.Draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screen.Layout(outsideWidth, outsideHeight)
}

func (g *game) startDailyLevel() {
	today12am := time.Now().Truncate(24 * time.Hour)
	g.screen = NewLevelScreen(today12am.UnixNano(), g.backToMainMenu)
}

func (g *game) startRandomLevel() {
	seed := time.Now().UTC().UnixNano()
	g.screen = NewLevelScreen(seed, g.backToMainMenu)
}

func (g *game) backToMainMenu() {
	g.screen = NewMainMenuScreen(
		g.startDailyLevel,
		g.startRandomLevel,
	)
}
