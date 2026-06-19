package collision

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/partitioning/hashgrid"
	"github.com/yohamta/donburi"
)

type CollisionPlugin struct {
	staticEntities  *hashgrid.HashGrid[donburi.Entity]
	dynamicEntities *hashgrid.HashGrid[donburi.Entity]
}

func NewCollisionPlugin(world engine.World) *CollisionPlugin {
	return &CollisionPlugin{
		staticEntities:  hashgrid.New[donburi.Entity](32, hashgrid.Padding{}),
		dynamicEntities: hashgrid.New[donburi.Entity](32, hashgrid.Padding{}),
	}
}
