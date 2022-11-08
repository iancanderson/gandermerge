package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/layers"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type render struct {
	query      *query.Query
	scoreQuery *query.Query
}

var Render = &render{
	query: ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.Sprite,
			component.Selectable,
		)),
	scoreQuery: ecs.NewQuery(
		layers.LayerScoreboard,
		filter.Contains(
			component.Score,
		)),
}

func (r *render) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	// Check if game is over
	score, ok := r.scoreQuery.FirstEntity(ecs.World)
	gameOverColorScale := 1.0
	if ok && component.Score.Get(score).GameOver() {
		gameOverColorScale = 0.25
	}

	r.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.Sprite.Get(entry)
		selectable := component.Selectable.Get(entry)
		op := sprite.DrawOptions()
		if selectable.Selected {
			op.ColorM.Scale(0.5, 0.5, 0.5, 1)
		}
		op.ColorM.Scale(1.0, 1.0, 1.0, gameOverColorScale)

		screen.DrawImage(sprite.Image, op)
	})
}
