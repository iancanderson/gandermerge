package component

import "github.com/yohamta/donburi"

type ProjectileData struct {
	DestX float64
	DestY float64
}

var Projectile = donburi.NewComponentType[ProjectileData]()

func GetProjectile(entry *donburi.Entry) *ProjectileData {
	return donburi.Get[ProjectileData](entry, Projectile)
}
