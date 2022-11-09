package component

import "github.com/yohamta/donburi"

type ModalData struct {
	Active bool
	Text   string
}

var Modal = donburi.NewComponentType[ModalData]()
