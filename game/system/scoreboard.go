package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/font/inconsolata"
)

type scoreboard struct {
	query *query.Query
}

var Scoreboard = &scoreboard{
	query: ecs.NewQuery(
		layers.LayerScoreboard,
		filter.Contains(
			component.Score,
		)),
}

func (s *scoreboard) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	entry, ok := s.query.FirstEntity(ecs.World)
	if !ok {
		return
	}

	score := component.GetScore(entry)
	moves := fmt.Sprintf("Moves Remaining: %d", score.MovesRemaining)
	text.Draw(screen, moves, inconsolata.Bold8x16, 20, 30, color.White)

	if score.TotalEnergyGoal <= 0 {
		text.Draw(screen, "You Win!", inconsolata.Bold8x16, 20, 50, color.RGBA{0x00, 0xff, 0x00, 0xff})
	} else {
		energyToWin := fmt.Sprintf("Energy to Win: %d", score.TotalEnergyGoal)
		text.Draw(screen, energyToWin, inconsolata.Bold8x16, 20, 50, color.White)
	}
}
