package component

import "github.com/yohamta/donburi"

type SelectableData struct {
	Selected bool
}

var Selectable = donburi.NewComponentType[SelectableData]()

func GetSelectable(entry *donburi.Entry) *SelectableData {
	return donburi.Get[SelectableData](entry, Selectable)
}
