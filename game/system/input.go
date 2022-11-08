package system

import (
	"image/color"
	"time"

	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/core"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/iancanderson/gandermerge/game/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type orbChain struct {
	orbs       []*donburi.Entry
	energyType core.EnergyType
}

func (o *orbChain) IsSecondToLast(entry *donburi.Entry) bool {
	return len(o.orbs) > 1 && o.orbs[len(o.orbs)-2] == entry
}

func (o *orbChain) Add(entry *donburi.Entry) bool {
	if o.contains(entry) {
		return false
	}
	if o.energyType != component.GetEnergy(entry).EnergyType {
		return false
	}
	if len(o.orbs) > 0 {
		// Only allow adjacent orbs to be added to the chain
		lastOrb := o.orbs[len(o.orbs)-1]
		if !component.GetGridPosition(lastOrb).IsAdjacent(component.GetGridPosition(entry)) {
			return false
		}
	}

	o.orbs = append(o.orbs, entry)
	return true
}

func (o *orbChain) CanBeMerged() bool {
	return len(o.orbs) >= 3
}

func (o *orbChain) Pop() *donburi.Entry {
	// Return the last orb in the chain and remove it from the chain
	i := len(o.orbs) - 1
	orb := o.orbs[i]
	o.orbs = o.orbs[:i]
	return orb
}

func (o *orbChain) Len() int {
	return len(o.orbs)
}

func (o *orbChain) contains(entry *donburi.Entry) bool {
	for _, orb := range o.orbs {
		if orb == entry {
			return true
		}
	}
	return false
}

type input struct {
	chain           *orbChain
	inputSource     util.InputSource
	scoreQuery      *query.Query
	selectableQuery *query.Query
	enemyQuery      *query.Query
	soundManager    *util.SoundManager
}

var Input = &input{
	enemyQuery: ecs.NewQuery(
		layers.LayerEnemy,
		filter.Contains(
			component.Energy,
			component.Hitpoints,
		)),
	selectableQuery: ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.Energy,
			component.Selectable,
			component.Sprite,
			component.GridPosition,
		)),
	scoreQuery: ecs.NewQuery(
		layers.LayerScoreboard,
		filter.Contains(
			component.Score,
		)),
	soundManager: util.NewSoundManager(),
}

func (r *input) Startup(ecs *ecs.ECS) {
	r.soundManager.LoadSounds()
}

func (r *input) Update(ecs *ecs.ECS) {
	// Check if game is over
	score, ok := r.scoreQuery.FirstEntity(ecs.World)
	if ok && component.GetScore(score).IsGameOver() {
		return
	}

	if r.chain == nil {
		r.inputSource = util.JustPressedInputSource()
		if r.inputSource != nil {
			r.handleInputPressed(ecs)
		}
	}

	if r.inputSource != nil && r.inputSource.JustReleased() {
		r.clearOrbChain(ecs, ecs.World)
	}

	if r.chain != nil {
		r.selectableQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
			sprite := component.GetSprite(entry)
			inputX, inputY := r.inputSource.Position()
			if sprite.In(inputX, inputY) {
				if r.chain.Add(entry) {
					selectable := component.GetSelectable(entry)
					selectable.Selected = true
				} else if r.chain.IsSecondToLast(entry) {
					popped := r.chain.Pop()
					selectable := component.GetSelectable(popped)
					selectable.Selected = false
				}
			}
		})
	}
}

func (r *input) handleInputPressed(ecs *ecs.ECS) {
	r.selectableQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.GetSprite(entry)
		inputX, inputY := r.inputSource.Position()
		if sprite.In(inputX, inputY) {
			r.chain = r.createOrbChain(entry)
			selectable := component.GetSelectable(entry)
			selectable.Selected = true
			r.soundManager.PlayChainSound(r.chain.energyType)
		}
	})
}

func (r *input) createOrbChain(entry *donburi.Entry) *orbChain {
	chain := &orbChain{
		orbs:       []*donburi.Entry{entry},
		energyType: component.GetEnergy(entry).EnergyType,
	}
	return chain
}

func (r *input) clearOrbChain(ecs *ecs.ECS, world donburi.World) {
	if r.chain == nil {
		return
	}

	if r.chain.CanBeMerged() {
		r.soundManager.PauseChainSound(r.chain.energyType)
		r.soundManager.PlayMergeSound(r.chain.energyType)
		r.hitEnemy(ecs, world)

		for _, orb := range r.chain.orbs {
			donburi.Add(orb, component.Projectile,
				&component.ProjectileData{
					DestX: config.WindowWidth/2 - config.ColumnWidth/2,
					DestY: 250,
				})
			orb.RemoveComponent(component.GridPosition)
		}
	} else {
		for _, orb := range r.chain.orbs {
			selectable := component.GetSelectable(orb)
			selectable.Selected = false
		}

		r.soundManager.PauseChainSound(r.chain.energyType)
	}

	r.chain = nil
}

func (r *input) hitEnemy(ecs *ecs.ECS, world donburi.World) {
	energyEmitted := r.chain.Len()
	entry, ok := r.scoreQuery.FirstEntity(world)
	if ok {
		score := component.GetScore(entry)
		score.MovesRemaining--

		enemyEntry, ok := r.enemyQuery.FirstEntity(world)
		if ok {
			enemyEnergyType := component.GetEnergy(enemyEntry).EnergyType
			attackStrength := core.ScaleAttack(energyEmitted, r.chain.energyType, enemyEnergyType)
			component.Hitpoints.Get(enemyEntry).Hitpoints -= attackStrength

			multiplier := core.AttackMultiplier(r.chain.energyType, enemyEnergyType)
			spawnMultiplierSign(ecs, world, multiplier)
		}
	}
}

func spawnMultiplierSign(ecs *ecs.ECS, world donburi.World, multiplier float64) {
	entity := ecs.Create(layers.LayerEnemy, component.Text, component.Expiration)
	entry := ecs.World.Entry(entity)

	multiplierStr := "1x"
	var multiplierColor color.Color = color.White
	if multiplier == 0.5 {
		multiplierStr = "½x"
		multiplierColor = color.RGBA{0xff, 0x00, 0x00, 0xff}
	} else if multiplier == 2 {
		multiplierStr = "2x"
		multiplierColor = color.RGBA{0x00, 0xff, 0x00, 0xff}
	}

	donburi.SetValue(entry, component.Text,
		component.TextData{
			Text:     multiplierStr,
			X:        100,
			Y:        300,
			FontFace: util.FontManager.Go108,
			Color:    multiplierColor,
		})

	donburi.SetValue(entry, component.Expiration,
		component.ExpirationData{
			TTL: time.Second,
		})
}
