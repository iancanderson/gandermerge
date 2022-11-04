package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/iancanderson/gandermerge/component"
	"github.com/iancanderson/gandermerge/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type orbChain struct {
	orbs       []*donburi.Entry
	energyType component.EnergyType
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

func (o *orbChain) contains(entry *donburi.Entry) bool {
	for _, orb := range o.orbs {
		if orb == entry {
			return true
		}
	}
	return false
}

type input struct {
	selectableQuery *query.Query
	chain           *orbChain
}

var Input = &input{
	selectableQuery: ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.Energy,
			component.Selectable,
			component.Sprite,
		)),
}

//TODO: see example for input system: https://github.com/hajimehoshi/ebiten/blob/main/examples/drag/main.go

type InputSource interface {
	Position() (int, int)
}

type MouseInputSource struct{}

func (m *MouseInputSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (r *input) Update(ecs *ecs.ECS) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		inputSource := MouseInputSource{}

		r.selectableQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
			sprite := component.GetSprite(entry)
			inputX, inputY := inputSource.Position()
			if sprite.In(inputX, inputY) {
				r.chain = r.createOrbChain(entry)
				selectable := component.GetSelectable(entry)
				selectable.Selected = true
			}
		})
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		r.clearOrbChain(ecs.World)
	}

	if r.chain != nil {
		inputSource := MouseInputSource{}

		r.selectableQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
			sprite := component.GetSprite(entry)
			inputX, inputY := inputSource.Position()
			if sprite.In(inputX, inputY) {
				if r.chain.Add(entry) {
					selectable := component.GetSelectable(entry)
					selectable.Selected = true
				}
			}
		})
	}
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
	for _, orb := range r.chain.orbs {
		if r.chain.CanBeMerged() {
			world.Remove(orb.Entity())
		} else {
			selectable := component.GetSelectable(orb)
			selectable.Selected = false
		}
	}
	r.chain = nil
}
