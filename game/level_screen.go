package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/spookypaths/game/assets"
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/iancanderson/spookypaths/game/system"
	"github.com/iancanderson/spookypaths/game/uicomponent"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/furex/v2"
)

type LevelScreen struct {
	backToMainMenu func()
	ecs            *ecs.ECS
	gameUI         *furex.View
	gameOverText   *TextComponent
}

var scoreQuery = ecs.NewQuery(
	layers.LayerUI,
	filter.Contains(
		component.Score,
	))

func (g *LevelScreen) Update() error {
	g.ecs.Update()

	score, ok := scoreQuery.FirstEntity(g.ecs.World)
	if ok && component.Score.Get(score).GameOver() {
		if component.Score.Get(score).Won() {
			g.gameOverText.Str = "You won!"
		} else {
			g.gameOverText.Str = "You lost!"
		}
		g.gameUI.Update()
	}
	return nil
}

func (g *LevelScreen) Draw(screen *ebiten.Image) {
	screen.Clear()
	for layer := layers.LayerBackground; layer <= layers.LayerMetrics; layer++ {
		g.ecs.DrawLayer(layer, screen)
	}

	score, ok := scoreQuery.FirstEntity(g.ecs.World)
	if ok && component.Score.Get(score).GameOver() {
		g.gameUI.Draw(screen)
	}
}

func (g *LevelScreen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.WindowWidth, config.WindowHeight
}

func NewLevelScreen(seed int64, backToMainMenu func()) *LevelScreen {
	rand.Seed(seed)

	g := &LevelScreen{
		ecs:            createECS(),
		backToMainMenu: backToMainMenu,
	}

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
	g.ecs.AddRenderer(layers.LayerUI, system.Scoreboard.Draw)
	g.ecs.AddSystem(system.Expirator.Update)
	g.ecs.AddSystem(system.Modal.Update)
	g.ecs.AddRenderer(layers.LayerModal, system.Modal.Draw)

	g.setupMenuUI()

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

func (g *LevelScreen) setupMenuUI() {
	g.gameUI = &furex.View{
		Width:        config.WindowWidth,
		Height:       config.WindowHeight,
		Direction:    furex.Column,
		Justify:      furex.JustifyCenter,
		AlignItems:   furex.AlignItemCenter,
		AlignContent: furex.AlignContentCenter,
		Wrap:         furex.Wrap,
	}

	g.gameOverText = &TextComponent{fontface: assets.FontManager.Creepster160, Str: ""}
	g.gameUI.AddChild(&furex.View{
		Width:        config.WindowWidth - 100,
		Height:       120,
		Handler:      g.gameOverText,
		MarginBottom: 100,
		MarginTop:    200,
	})

	g.gameUI.AddChild(&furex.View{
		MarginBottom: 20,
		Width:        600,
		Height:       100,
		Handler: &uicomponent.Button{
			Text:    "Main menu",
			OnClick: g.backToMainMenu,
		},
	})
}
