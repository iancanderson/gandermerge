package component

import "github.com/yohamta/donburi"

type ProjectileData struct {
	DestX float64
	DestY float64
}

var Projectile = donburi.NewComponentType[ProjectileData]()
