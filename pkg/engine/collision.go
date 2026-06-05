package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/partitioning/spatialhash"
	"github.com/yohamta/donburi"
)

type Collision interface {
	Add(entry *donburi.Entry) bool
	Remove(entry *donburi.Entry) bool

	EnableCollision(a, b colliders.CollisionLayer)
	DisableCollision(a, b colliders.CollisionLayer)
	CanCollide(a, b colliders.CollisionLayer) bool
}

type (
	collisionIndexing map[donburi.Entity]spatialhash.SpatialIndex
	collisionMask     map[colliders.CollisionLayer]colliders.CollisionLayer
)

type collision struct {
	static  *spatialhash.SpatialHash[donburi.Entity]
	dynamic *spatialhash.SpatialHash[donburi.Entity]

	indexing collisionIndexing
	masks    collisionMask
}

func newCollision() *collision {
	return &collision{
		indexing: make(collisionIndexing),
		masks:    make(collisionMask),
	}
}

func (c *collision) Add(entry *donburi.Entry) bool {
	return true
}

func (c *collision) Remove(entry *donburi.Entry) bool {
	return true
}

func (c *collision) EnableCollision(a, b colliders.CollisionLayer) {
	c.masks[a] |= b
	c.masks[b] |= a
}

func (c *collision) DisableCollision(a, b colliders.CollisionLayer) {
	c.masks[a] &^= b
	c.masks[b] &^= a
}

func (c *collision) CanCollide(a, b colliders.CollisionLayer) bool {
	if a == b {
		return true
	}
	return (c.masks[a] & b) != 0
}

func (c *collision) checkCollision(world donburi.World) error {
	return nil
}
