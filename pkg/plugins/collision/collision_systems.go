package collision

import (
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type CollisionInfo struct {
}

type CollisionSystem struct {
	world *CollisionWorld
	query *donburi.Query

	collisions   []CollisionInfo
	collisionIdx int
}

func NewCollisionSystem(world *CollisionWorld) *CollisionSystem {
	return &CollisionSystem{
		world: world,
		query: donburi.NewQuery(
			filter.Contains(
				Collision,
				Collider,
				transform.Transform,
			),
		),
		collisions: make([]CollisionInfo, 0),
	}
}

func (c *CollisionSystem) CheckStaticCollision(ecs donburi.World, region geom.AABB) ([]CollisionInfo, bool) {
	c.collisions = c.collisions[:0]
	c.collisionIdx = 0

	return c.collisions, len(c.collisions) > 0
}
