package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/iancanderson/gandermerge/game/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const WindowHeight = 384
const WindowWidth = 384

type Game struct {
	ecs *ecs.ECS
}

func (g *Game) Update() error {
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ecs.DrawLayer(layers.LayerBackground, screen)
	g.ecs.DrawLayer(layers.LayerOrbs, screen)
	g.ecs.DrawLayer(layers.LayerMetrics, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func NewGame() *Game {
	g := &Game{
		ecs: createECS(),
	}
	orbSpawner := system.NewOrbSpawner()
	orbSpawner.Startup(g.ecs)

	g.ecs.AddSystems(
		ecs.System{
			Update: system.Input.Update,
		},
		ecs.System{
			Update: system.GridGravity.Update,
		},
		ecs.System{
			Update: orbSpawner.Update,
		},
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
