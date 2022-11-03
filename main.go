package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/layers"
	"github.com/iancanderson/gandermerge/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const windowHeight = 384
const windowWidth = 384

type Game struct {
	ecs *ecs.ECS
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ecs.DrawLayer(layers.LayerBackground, screen)
	g.ecs.DrawLayer(layers.LayerOrbs, screen)
	g.ecs.DrawLayer(layers.LayerMetrics, screen)
	// ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func NewGame() *Game {
	g := &Game{
		ecs: createECS(),
	}
	orbSpawner := system.NewOrbSpawner()
	orbSpawner.Startup(g.ecs)

	g.ecs.AddSystems(
		// ecs.System{
		// 	Update: orbSpawner.Update,
		// },
		ecs.System{
			Layer: layers.LayerOrbs,
			Draw:  system.Render.Draw,
		},
	)
	return g
}

func createECS() *ecs.ECS {
	world := createWorld()
	ecs := ecs.NewECS(world)
	return ecs
}

func createWorld() donburi.World {
	world := donburi.NewWorld()
	return world
}

func main() {
	ebiten.SetWindowTitle("Gandermerge")

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowSizeLimits(300, 200, -1, -1)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	rand.Seed(time.Now().UTC().UnixNano())
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
