package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/core"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type enemy struct {
	hitpointsBar hitpointsBar
	images       map[core.EnergyType]*ebiten.Image
	scoreQuery   *query.Query
	sprites      *query.Query
	textQuery    *query.Query
}

var Enemy = &enemy{
	hitpointsBar: hitpointsBar{
		hpMax: config.EnemyHitpoints,
		hp:    config.EnemyHitpoints,
	},
	sprites: ecs.NewQuery(
		layers.LayerEnemy,
		filter.Contains(
			component.Sprite,
		)),
	scoreQuery: ecs.NewQuery(
		layers.LayerScoreboard,
		filter.Contains(
			component.Score,
		)),
	textQuery: ecs.NewQuery(
		layers.LayerEnemy,
		filter.Contains(
			component.Text,
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
}

func (e *enemy) Startup(ecs *ecs.ECS) {
	//TODO: share these with orb_spawner?
	e.images = loadEnergyTypeImages()
	e.spawnEnemy(ecs)
}

func (e *enemy) NewGame(ecs *ecs.ECS) {
	e.sprites.EachEntity(ecs.World, func(entry *donburi.Entry) {
		ecs.World.Remove(entry.Entity())
	})
	e.spawnEnemy(ecs)
}

func (e *enemy) spawnEnemy(ecs *ecs.ECS) {
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
	score := component.Score.Get(scoreEntry)
	if score.Won() {
		e.sprites.EachEntity(ecs.World, func(entry *donburi.Entry) {
			ecs.World.Remove(entry.Entity())
		})

		HitpointsBar.hide = true
	}
}

func (e *enemy) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	// TODO: consolidate with render.go
	e.sprites.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.Sprite.Get(entry)
		op := sprite.DrawOptions()
		screen.DrawImage(sprite.Image, op)
	})

	e.textQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
		textEntry := component.Text.Get(entry)
		if textEntry.Bubble == component.BubbleLeft {
			textRect := text.BoundString(textEntry.FontFace, textEntry.Text)
			radius := float64(textRect.Dx())
			ebitenutil.DrawCircle(screen, float64(textEntry.X+textRect.Dx()/2), float64(textEntry.Y)-float64(textRect.Dy()/2), radius, color.White)
			ebitenutil.DrawCircle(screen, float64(textEntry.X+textRect.Dx()/2)-100, float64(textEntry.Y)-float64(textRect.Dy()/2)-70, radius/5, color.White)
			ebitenutil.DrawCircle(screen, float64(textEntry.X+textRect.Dx()/2)-135, float64(textEntry.Y)-float64(textRect.Dy()/2)-40, radius/7, color.White)
		}
		text.Draw(screen, textEntry.Text, textEntry.FontFace, textEntry.X, textEntry.Y, textEntry.Color)
	})

	HitpointsBar.Draw(ecs, screen)
}
