package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/spookypaths/game/assets/images"
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/iancanderson/spookypaths/game/util"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type modal struct {
	infoButtonQuery   *query.Query
	infoButtonPressed bool
	inputSource       util.InputSource
	modalQuery        *query.Query
}

var Modal = &modal{
	infoButtonQuery: ecs.NewQuery(
		layers.LayerModal,
		filter.Contains(
			component.Sprite,
			component.InfoButton,
		)),
	modalQuery: ecs.NewQuery(
		layers.LayerModal,
		filter.Contains(
			component.Modal,
		)),
}

func (m *modal) Startup(ecs *ecs.ECS) {
	m.spawnInfoButton(ecs)

	modal := ecs.Create(
		layers.LayerModal,
		component.Modal,
	)
	entry := ecs.World.Entry(modal)
	component.Modal.SetValue(entry, component.ModalData{
		Active: false,
		Text:   config.ModalText,
	})
}

func (m *modal) spawnInfoButton(ecs *ecs.ECS) {
	info := ecs.Create(
		layers.LayerModal,
		component.Sprite,
		component.InfoButton,
	)
	entry := ecs.World.Entry(info)
	component.Sprite.Set(entry, component.NewSpriteData(
		util.LoadImage(images.Information_png),
		config.WindowWidth-100,
		10,
	).WithScale(0.15))
}

func (m *modal) Update(ecs *ecs.ECS) {
	infoButton, ok := m.infoButtonQuery.FirstEntity(ecs.World)
	if !ok {
		panic("no info button")
	}
	buttonSprite := component.Sprite.Get(infoButton)
	if m.inputSource == nil {
		m.inputSource = util.JustPressedInputSource()
	}
	if m.inputSource != nil {
		inputX, inputY := m.inputSource.Position()
		if buttonSprite.CloseTo(inputX, inputY) {
			m.infoButtonPressed = true
		}

		if m.infoButtonPressed && m.inputSource.JustReleased() {
			m.infoButtonPressed = false
			m.inputSource = nil

			modal, ok := m.modalQuery.FirstEntity(ecs.World)
			if !ok {
				panic("no modal")
			}
			modalEntry := component.Modal.Get(modal)
			modalEntry.Active = !modalEntry.Active
		}
	}
}

func (m *modal) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	modal, ok := m.modalQuery.FirstEntity(ecs.World)
	if !ok {
		panic("no modal")
	}
	modalEntry := component.Modal.Get(modal)
	if modalEntry.Active {
		bg := ebiten.NewImage(config.WindowWidth, config.WindowHeight)
		bg.Fill(color.White)
		screen.DrawImage(bg, nil)
		text.Draw(screen, modalEntry.Text, util.FontManager.Go36, 40, 70, color.Black)
	}
	m.drawInfoButton(ecs, screen)
}

func (m *modal) drawInfoButton(ecs *ecs.ECS, screen *ebiten.Image) {
	entry, ok := m.infoButtonQuery.FirstEntity(ecs.World)
	if !ok {
		panic("no info button")
	}
	sprite := component.Sprite.Get(entry)
	op := sprite.DrawOptions()
	screen.DrawImage(sprite.Image, op)
}
