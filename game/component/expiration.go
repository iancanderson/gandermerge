package component

import (
	"time"

	"github.com/yohamta/donburi"
)

type ExpirationData struct {
	TTL time.Duration
}

var Expiration = donburi.NewComponentType[ExpirationData]()
