package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/spookypaths/game/assets"
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/iancanderson/spookypaths/game/uicomponent"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"github.com/yohamta/furex/v2"
)

var scoreQuery = ecs.NewQuery(
	layers.LayerUI,
	filter.Contains(
		component.Score,
	))

type levelUI struct {
	modalQuery *query.Query
	modalView  *furex.View
	view       *furex.View
	modalBg    *ebiten.Image
}

var LevelUI = &levelUI{
	modalBg: ebiten.NewImage(config.WindowWidth, config.WindowHeight),
	modalQuery: ecs.NewQuery(
		layers.LayerModal,
		filter.Contains(
			component.Modal,
		)),
}

func (ui *levelUI) Startup(ecs *ecs.ECS) {
	ui.modalBg.Fill(color.White)

	modal := ecs.Create(
		layers.LayerModal,
		component.Modal,
	)
	entry := ecs.World.Entry(modal)
	component.Modal.SetValue(entry, component.ModalData{
		Active: false,
		Text:   config.ModalText,
	})

	ui.view = (&furex.View{
		Width:     config.WindowWidth,
		Height:    config.WindowHeight,
		Direction: furex.Row,
		Justify:   furex.JustifySpaceBetween,
	}).AddChild(&furex.View{
		Width:      80,
		Height:     80,
		MarginTop:  20,
		MarginLeft: 20,
		Handler: &uicomponent.Toggle{
			ImageOn:  assets.ImageManager.SoundOn,
			ImageOff: assets.ImageManager.SoundOff,
			On:       true,
			OnToggle: assets.SoundManager.Toggle,
		},
	}).AddChild(&furex.View{
		Width:       80,
		Height:      80,
		MarginTop:   20,
		MarginRight: 20,
		Handler: &uicomponent.Button{
			Text: "?",
			OnClick: func() {
				modal, ok := ui.modalQuery.FirstEntity(ecs.World)
				if !ok {
					panic("no modal")
				}
				modalEntry := component.Modal.Get(modal)
				modalEntry.Active = !modalEntry.Active
			},
		},
	})

	ui.modalView = (&furex.View{
		Width:     config.WindowWidth,
		Height:    config.WindowHeight,
		Direction: furex.Row,
		Justify:   furex.JustifyEnd,
	}).AddChild(&furex.View{
		Width:       80,
		Height:      80,
		MarginTop:   20,
		MarginRight: 20,
		Handler: &uicomponent.Button{
			Text: "?",
			OnClick: func() {
				modal, ok := ui.modalQuery.FirstEntity(ecs.World)
				if !ok {
					panic("no modal")
				}
				modalEntry := component.Modal.Get(modal)
				modalEntry.Active = !modalEntry.Active
			},
		},
	})
}

func (ui *levelUI) Update(ecs *ecs.ECS) {
	modal, ok := ui.modalQuery.FirstEntity(ecs.World)
	if !ok {
		panic("no modal")
	}
	modalEntry := component.Modal.Get(modal)

	if modalEntry.Active {
		ui.modalView.Update()
	} else {
		ui.view.Update()
	}
}

func (ui *levelUI) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	ui.view.Draw(screen)
}

func (ui *levelUI) DrawModal(ecs *ecs.ECS, screen *ebiten.Image) {
	score, ok := scoreQuery.FirstEntity(ecs.World)
	if ok && component.Score.Get(score).GameOver() {
		return
	}

	modal, ok := ui.modalQuery.FirstEntity(ecs.World)
	if !ok {
		panic("no modal")
	}
	modalEntry := component.Modal.Get(modal)
	if modalEntry.Active {
		screen.DrawImage(ui.modalBg, nil)
		text.Draw(screen, "Spooky Paths", assets.FontManager.Creepster72, 40, 100, color.Black)
		text.Draw(screen, modalEntry.Text, assets.FontManager.Mona36, 40, 180, color.Black)
		ui.modalView.Draw(screen)
	}
}
