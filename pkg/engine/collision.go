package engine

import (
	"github.com/adm87/onyx-game/pkg/engine/partitions"
	"github.com/yohamta/donburi"
)

type Collision interface {
}

type collision struct {
	partitions partitions.SpatialHash[donburi.Entity]
}

func newCollision() *collision {
	return &collision{}
}
