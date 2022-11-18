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

type modal struct {
	modalQuery *query.Query
	modalUI    *furex.View
}

var Modal = &modal{
	modalQuery: ecs.NewQuery(
		layers.LayerModal,
		filter.Contains(
			component.Modal,
		)),
}

func (m *modal) Startup(ecs *ecs.ECS) {
	modal := ecs.Create(
		layers.LayerModal,
		component.Modal,
	)
	entry := ecs.World.Entry(modal)
	component.Modal.SetValue(entry, component.ModalData{
		Active: false,
		Text:   config.ModalText,
	})

	m.modalUI = (&furex.View{
		Width:      config.WindowWidth,
		Height:     config.WindowHeight,
		AlignItems: furex.AlignItemEnd,
		Direction:  furex.Column,
	}).AddChild(&furex.View{
		Width:       80,
		Height:      80,
		MarginTop:   20,
		MarginRight: 20,
		Handler: &uicomponent.Button{
			Text: "?",
			OnClick: func() {
				modal, ok := m.modalQuery.FirstEntity(ecs.World)
				if !ok {
					panic("no modal")
				}
				modalEntry := component.Modal.Get(modal)
				modalEntry.Active = !modalEntry.Active
			},
		},
	})
}

func (m *modal) Update(ecs *ecs.ECS) {
	m.modalUI.Update()
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
		text.Draw(screen, modalEntry.Text, assets.FontManager.Mona36, 40, 70, color.Black)
	}

	m.modalUI.Draw(screen)
}
