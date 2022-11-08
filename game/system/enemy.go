package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/core"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type enemy struct {
	images       map[core.EnergyType]*ebiten.Image
	query        *query.Query
	scoreQuery   *query.Query
	hitpointsBar hitpointsBar
}

var Enemy = &enemy{
	hitpointsBar: hitpointsBar{
		hpMax: config.EnemyHitpoints,
		hp:    config.EnemyHitpoints,
	},
	query: ecs.NewQuery(
		layers.LayerEnemy,
		filter.Contains(
			component.Sprite,
		)),
	scoreQuery: ecs.NewQuery(
		layers.LayerScoreboard,
		filter.Contains(
			component.Score,
		)),
}

const enemyWidth = 309

type hitpointsBar struct {
	hpMax  int
	hp     int
	query  *query.Query
	width  int
	height int
	y      int
	hide   bool
}

var HitpointsBar = hitpointsBar{
	hpMax:  config.EnemyHitpoints,
	height: 20,
	width:  100,
	y:      400,
	hide:   false,
	query: ecs.NewQuery(
		layers.LayerEnemy,
		filter.Contains(
			component.Hitpoints,
		)),
}

func (h *hitpointsBar) Update(ecs *ecs.ECS) {
	enemy, ok := h.query.FirstEntity(ecs.World)
	if !ok {
		return
	}
	h.hp = component.Hitpoints.Get(enemy).Hitpoints
}

func (h *hitpointsBar) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	if h.hide {
		return
	}

	// outer rectangle
	ebitenutil.DrawRect(
		screen,
		config.WindowWidth/2-float64(h.width/2),
		float64(h.y),
		float64(h.width),
		float64(h.height),
		color.RGBA{0x00, 0xff, 0x00, 0xff},
	)

	widthScale := float64(h.hp) / float64(h.hpMax)
	innerX := config.WindowWidth/2 - float64(h.width/2) + 2
	innerWidth := float64(h.width) - 4
	progressWidth := widthScale * innerWidth
	ebitenutil.DrawRect(
		screen,
		innerX,
		float64(h.y+2),
		progressWidth,
		float64(h.height-4),
		color.RGBA{0xff, 0x00, 0x00, 0xff},
	)

	ebitenutil.DrawRect(
		screen,
		progressWidth+innerX,
		float64(h.y+2),
		innerWidth-progressWidth,
		float64(h.height-4),
		color.Black,
	)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Enemy HP: %d", h.hp))
}

func (e *enemy) Startup(ecs *ecs.ECS) {
	//TODO: share these with orb_spawner?
	e.images = loadEnergyTypeImages()

	entity := ecs.Create(
		layers.LayerEnemy,
		component.Energy,
		component.Sprite,
		component.Hitpoints,
	)
	entry := ecs.World.Entry(entity)

	energyType := component.RandomEnergyType()
	donburi.SetValue(entry, component.Energy,
		component.EnergyData{
			EnergyType: energyType,
		})

	donburi.SetValue(entry, component.Sprite,
		component.SpriteData{
			Image: e.images[energyType],
			X:     config.WindowWidth/2 - enemyWidth/2,
			Y:     100,
		}.WithScale(0.5).WithGreenTint(energyType == core.Poison).WithRedTint(energyType == core.Fire))

	donburi.SetValue(entry, component.Hitpoints,
		component.HitpointsData{
			MaxHitpoints: config.EnemyHitpoints,
			Hitpoints:    config.EnemyHitpoints,
		})

	HitpointsBar.hide = false
}

func (e *enemy) Update(ecs *ecs.ECS) {
	HitpointsBar.Update(ecs)

	scoreEntry, ok := e.scoreQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}
	score := component.GetScore(scoreEntry)
	if score.Won() {
		e.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
			ecs.World.Remove(entry.Entity())
		})

		HitpointsBar.hide = true
	}
}

func (e *enemy) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	// TODO: consolidate with render.go
	e.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.GetSprite(entry)
		op := sprite.DrawOptions()
		screen.DrawImage(sprite.Image, op)
	})

	HitpointsBar.Draw(ecs, screen)
}
