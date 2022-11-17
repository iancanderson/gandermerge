package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/spookypaths/game/assets"
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type scoreboard struct {
	scoreQuery *query.Query
}

var Scoreboard = &scoreboard{
	scoreQuery: ecs.NewQuery(
		layers.LayerUI,
		filter.Contains(
			component.Score,
		)),
}

func (s *scoreboard) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	scoreEntry, ok := s.scoreQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}

	score := component.Score.Get(scoreEntry)

	fontface := assets.FontManager.Go36
	var movesLeftColor color.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
	if score.MovesRemaining <= 5 {
		movesLeftColor = color.RGBA{0xff, 0x00, 0x00, 0xff}
	} else if score.MovesRemaining <= 10 {
		movesLeftColor = color.RGBA{0xff, 0xff, 0x00, 0xff}
	} else if score.MovesRemaining <= 15 {
		movesLeftColor = color.White
	}
	movePluralized := "moves"
	if score.MovesRemaining == 1 {
		movePluralized = "move"
	}
	moves := fmt.Sprintf("%d %s left", score.MovesRemaining, movePluralized)
	textWidth := text.BoundString(fontface, moves).Dx()
	text.Draw(screen, moves, fontface, config.WindowWidth/2-textWidth/2, 60, movesLeftColor)
}
