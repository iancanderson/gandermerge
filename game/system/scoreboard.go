package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/spookypaths/game/assets"
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/iancanderson/spookypaths/game/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type scoreboard struct {
	inputSource      util.InputSource
	scoreQuery       *query.Query
	playButtonQuery  *query.Query
	playAgainPressed bool
}

const buttonHeight = 60.0
const buttonWidth = 200.0

var Scoreboard = &scoreboard{
	playAgainPressed: false,
	playButtonQuery: ecs.NewQuery(
		layers.LayerUI,
		filter.Contains(component.Sprite),
	),
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

	if score.Won() {
		textValue := "You won!"
		textWidth = text.BoundString(fontface, textValue).Dx()
		text.Draw(screen, textValue, fontface, config.WindowWidth/2-textWidth/2, 110, color.RGBA{0x00, 0xff, 0x00, 0xff})
		s.drawPlayAgainButton(ecs, screen)
	} else if score.Lost() {
		textValue := "You lost!"
		textWidth = text.BoundString(fontface, textValue).Dx()
		text.Draw(screen, textValue, fontface, config.WindowWidth/2-textWidth/2, 110, color.RGBA{0xff, 0x00, 0x00, 0xff})
		s.drawPlayAgainButton(ecs, screen)
	}
}

func (s *scoreboard) Update(ecs *ecs.ECS) {
	scoreEntry, ok := s.scoreQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}

	score := component.Score.Get(scoreEntry)
	if score.GameOver() {
		playButton := s.findOrSpawnPlayAgainButton(ecs)
		buttonSprite := component.Sprite.Get(playButton)

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
				Enemy.NewGame(ecs)
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
	entity := ecs.Create(layers.LayerUI, component.Sprite)
	entry := ecs.World.Entry(entity)

	component.Sprite.Set(entry, component.NewSpriteData(
		ebiten.NewImage(buttonWidth, buttonHeight),
		config.WindowWidth/2-buttonWidth/2,
		config.WindowHeight/2,
	))

	return entry
}

func (s *scoreboard) drawPlayAgainButton(ecs *ecs.ECS, screen *ebiten.Image) {
	playButtonEntry, ok := s.playButtonQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}
	sprite := component.Sprite.Get(playButtonEntry)
	ebitenutil.DrawRect(sprite.Image, 0, 0, buttonWidth, buttonHeight, color.RGBA{0x00, 0xff, 0x00, 0xff})
	text.BoundString(assets.FontManager.Go36, "Play Again")
	text.Draw(sprite.Image, "Play again", assets.FontManager.Go36, 10, 40, color.Black)

	op := sprite.DrawOptions()
	screen.DrawImage(sprite.Image, op)
}
