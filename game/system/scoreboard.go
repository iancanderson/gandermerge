package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/iancanderson/gandermerge/game/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/font/inconsolata"
)

type scoreboard struct {
	inputSource      util.InputSource
	scoreQuery       *query.Query
	playButtonQuery  *query.Query
	playAgainPressed bool
}

const buttonHeight = 30.0
const buttonWidth = 100.0

var Scoreboard = &scoreboard{
	playAgainPressed: false,
	playButtonQuery: ecs.NewQuery(
		layers.LayerScoreboard,
		filter.Contains(component.Sprite),
	),
	scoreQuery: ecs.NewQuery(
		layers.LayerScoreboard,
		filter.Contains(
			component.Score,
		)),
}

func (s *scoreboard) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	scoreEntry, ok := s.scoreQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}

	score := component.GetScore(scoreEntry)
	moves := fmt.Sprintf("Moves Remaining: %d", score.MovesRemaining)
	text.Draw(screen, moves, inconsolata.Bold8x16, 20, 30, color.White)

	if score.Won() {
		text.Draw(screen, "You Win!", inconsolata.Bold8x16, 20, 50, color.RGBA{0x00, 0xff, 0x00, 0xff})
		s.drawPlayAgainButton(ecs, screen)
	} else if score.Lost() {
		text.Draw(screen, "You Lost!", inconsolata.Bold8x16, 20, 50, color.RGBA{0xff, 0x00, 0x00, 0xff})
		s.drawPlayAgainButton(ecs, screen)
	} else {
		energyToWin := fmt.Sprintf("Energy to Win: %d", score.EnergyToWin)
		text.Draw(screen, energyToWin, inconsolata.Bold8x16, 20, 50, color.White)
	}
}

func (s *scoreboard) Update(ecs *ecs.ECS) {
	scoreEntry, ok := s.scoreQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}

	score := component.GetScore(scoreEntry)
	if score.Won() || score.Lost() {
		playButton := s.findOrSpawnPlayAgainButton(ecs)
		buttonSprite := component.GetSprite(playButton)

		if s.inputSource == nil {
			s.inputSource = util.JustPressedInputSource()
		}
		if s.inputSource != nil {
			inputX, inputY := s.inputSource.Position()
			if buttonSprite.In(inputX, inputY) {
				s.playAgainPressed = true
			}

			if s.playAgainPressed && s.inputSource.JustReleased() {
				s.playAgainPressed = false
				s.inputSource = nil
				score.NewGame()
			}
		}
	}
}

func (s *scoreboard) findOrSpawnPlayAgainButton(ecs *ecs.ECS) *donburi.Entry {
	playButtonEntry, ok := s.playButtonQuery.FirstEntity(ecs.World)
	if !ok {
		playButtonEntry = spawnPlayAgainButton(ecs)
	}
	return playButtonEntry
}

func spawnPlayAgainButton(ecs *ecs.ECS) *donburi.Entry {
	entity := ecs.Create(layers.LayerScoreboard, component.Sprite)
	entry := ecs.World.Entry(entity)

	donburi.SetValue(entry, component.Sprite, component.SpriteData{
		Image: ebiten.NewImage(buttonWidth, buttonHeight),
		X:     config.WindowWidth/2 - buttonWidth/2,
		Y:     60,
		Scale: 1.0,
	})
	return entry
}

func (s *scoreboard) drawPlayAgainButton(ecs *ecs.ECS, screen *ebiten.Image) {
	playButtonEntry, ok := s.playButtonQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}
	sprite := component.GetSprite(playButtonEntry)
	ebitenutil.DrawRect(sprite.Image, 0, 0, buttonWidth, buttonHeight, color.RGBA{0x00, 0xff, 0x00, 0xff})
	text.Draw(sprite.Image, "Play Again", inconsolata.Bold8x16, 10, 20, color.Black)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sprite.Scale, sprite.Scale)
	op.GeoM.Translate(sprite.X+4, sprite.Y+4)

	screen.DrawImage(sprite.Image, op)
}
