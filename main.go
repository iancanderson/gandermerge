package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/game"
	"github.com/iancanderson/gandermerge/game/config"
)

func main() {
	ebiten.SetWindowTitle("Gandermerge")

	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowSizeLimits(300, 200, -1, -1)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	rand.Seed(time.Now().UTC().UnixNano())
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
