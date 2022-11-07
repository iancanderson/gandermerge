package system

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/iancanderson/gandermerge/game/assets/sounds"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/iancanderson/gandermerge/game/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type orbChain struct {
	orbs       []*donburi.Entry
	energyType component.EnergyType
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
	energySounds    map[component.EnergyType]*audio.Player
}

var Input = &input{
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
}

func (r *input) Startup(ecs *ecs.ECS) {
	audioContext := audio.NewContext(config.AudioSampleRate)

	soundData := map[component.EnergyType][]byte{
		component.Poison:  sounds.Poison_wav,
		component.Ghost:   sounds.Ghost_breath_wav,
		component.Psychic: sounds.Psychic_wav,
	}

	r.energySounds = make(map[component.EnergyType]*audio.Player)

	for energyType, data := range soundData {
		stream, err := wav.DecodeWithSampleRate(config.AudioSampleRate, bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		player, err := audioContext.NewPlayer(stream)
		if err != nil {
			panic(err)
		}
		r.energySounds[energyType] = player
	}
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
		r.clearOrbChain(ecs.World)
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

			sound := r.energySounds[r.chain.energyType]
			if sound != nil {
				sound.Rewind()
				sound.Play()
			}
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

func (r *input) clearOrbChain(world donburi.World) {
	if r.chain == nil {
		return
	}

	if r.chain.CanBeMerged() {
		energyEmitted := r.chain.Len()
		entry, ok := r.scoreQuery.FirstEntity(world)
		if ok {
			score := component.GetScore(entry)
			score.EnergyToWin -= energyEmitted
			score.MovesRemaining--
		}

		for _, orb := range r.chain.orbs {
			donburi.Add(orb, component.Projectile,
				&component.ProjectileData{
					DestX: config.WindowWidth/2 - config.ColumnWidth/2,
					DestY: 250,
				})
			orb.RemoveComponent(component.GridPosition)

			//TODO: remove them after they reach the enemy
			// world.Remove(orb.Entity())
		}
	} else {
		for _, orb := range r.chain.orbs {
			selectable := component.GetSelectable(orb)
			selectable.Selected = false
		}

		sound := r.energySounds[r.chain.energyType]
		if sound != nil {
			sound.Pause()
		}
	}

	r.chain = nil
}
