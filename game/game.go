package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/iancanderson/gandermerge/game/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Game struct {
	ecs *ecs.ECS
}

func (g *Game) Update() error {
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for layer := layers.LayerBackground; layer <= layers.LayerMetrics; layer++ {
		g.ecs.DrawLayer(layer, screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.WindowWidth, config.WindowHeight
}

func NewGame() *Game {
	g := &Game{
		ecs: createECS(),
	}
	orbSpawner := system.NewOrbSpawner()
	orbSpawner.Startup(g.ecs)

	scorer := system.Scorer
	scorer.Startup(g.ecs)

	system.Scoreboard.Startup(g.ecs)
	system.Enemy.Startup(g.ecs)
	system.Input.Startup(g.ecs)

	g.ecs.AddSystem(system.Input.Update)
	g.ecs.AddSystem(system.GridGravity.Update)
	g.ecs.AddSystem(orbSpawner.Update)
	g.ecs.AddRenderer(layers.LayerOrbs, system.Render.Draw)
	g.ecs.AddSystem(system.Projectile.Update)
	g.ecs.AddSystem(system.Enemy.Update)
	g.ecs.AddRenderer(layers.LayerEnemy, system.Enemy.Draw)
	g.ecs.AddSystem(system.Scoreboard.Update)
	g.ecs.AddRenderer(layers.LayerScoreboard, system.Scoreboard.Draw)

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
