package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/spookypaths/game/assets"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/iancanderson/spookypaths/game/system"
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
	rand.Seed(time.Now().UTC().UnixNano())

	g := &Game{
		ecs: createECS(),
	}

	loadAssets()

	orbSpawner := system.NewOrbSpawner()
	orbSpawner.Startup(g.ecs)

	scorer := system.Scorer
	scorer.Startup(g.ecs)

	system.Enemy.Startup(g.ecs)
	system.Modal.Startup(g.ecs)

	g.ecs.AddSystem(system.Input.Update)
	g.ecs.AddSystem(system.GridGravity.Update)
	g.ecs.AddSystem(orbSpawner.Update)
	g.ecs.AddRenderer(layers.LayerOrbs, system.Render.Draw)
	g.ecs.AddSystem(system.Projectile.Update)
	g.ecs.AddSystem(system.Enemy.Update)
	g.ecs.AddRenderer(layers.LayerEnemy, system.Enemy.Draw)
	g.ecs.AddSystem(system.Scorer.Update)
	g.ecs.AddSystem(system.Scoreboard.Update)
	g.ecs.AddRenderer(layers.LayerScoreboard, system.Scoreboard.Draw)
	g.ecs.AddSystem(system.Expirator.Update)
	g.ecs.AddSystem(system.Modal.Update)
	g.ecs.AddRenderer(layers.LayerModal, system.Modal.Draw)

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

func loadAssets() {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(3)

	var assetManagers []assets.Manager = []assets.Manager{
		assets.FontManager,
		assets.ImageManager,
		assets.SoundManager,
	}

	for _, assetManager := range assetManagers {
		go func(assetManager assets.Manager) {
			defer wg.Done()
			assetManager.Load()
		}(assetManager)
	}

	wg.Wait()
	fmt.Printf("Loaded assets in %v ms\n", time.Since(start).Milliseconds())
}
