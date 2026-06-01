package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning"
	"github.com/yohamta/donburi"
)

type Collision interface {
	Partitioning() *partitioning.SpatialHash[donburi.Entity]
	StaticPartitioning() *partitioning.SpatialHash[donburi.Entity]

	Add(entry *donburi.Entry) bool
	Remove(entry *donburi.Entry) bool
	Update(entry *donburi.Entry) bool
	Simulate()

	Query(aabb geom.AABB) []donburi.Entity
	QueryStatic(aabb geom.AABB) []donburi.Entity
}

type collision struct {
	staticPartitions *partitioning.SpatialHash[donburi.Entity]
	partitions       *partitioning.SpatialHash[donburi.Entity]

	entities   map[donburi.Entity]partitioning.SpatialIndex
	queryCache []donburi.Entity
}

func newCollision() *collision {
	return &collision{
		staticPartitions: partitioning.NewSpatialHash[donburi.Entity](100, 32, 64),
		partitions:       partitioning.NewSpatialHash[donburi.Entity](100, 16),
		entities:         make(map[donburi.Entity]partitioning.SpatialIndex),
	}
}

func (c *collision) StaticPartitioning() *partitioning.SpatialHash[donburi.Entity] {
	return c.staticPartitions
}

func (c *collision) Partitioning() *partitioning.SpatialHash[donburi.Entity] {
	return c.partitions
}

func (c *collision) Add(entry *donburi.Entry) bool {
	entity := entry.Entity()

	if _, exists := c.entities[entity]; exists {
		return false
	}

	collisionType := colliders.GetColliderType(entry)
	collider := colliders.GetBoxCollider(entry)
	position := transform.GetPosition(entry)

	var idx partitioning.SpatialIndex
	var ok bool

	switch collisionType {
	case colliders.ColliderTypeStatic:
		idx, ok = c.staticPartitions.Insert(entity, collider.Translate(position))
	default:
		idx, ok = c.partitions.Insert(entity, collider.Translate(position))
	}

	if !ok {
		return false
	}

	c.entities[entity] = idx
	return true
}

func (c *collision) Remove(entry *donburi.Entry) bool {
	entity := entry.Entity()

	idx, exists := c.entities[entity]
	if !exists {
		return false
	}

	collisionType := colliders.GetColliderType(entry)

	switch collisionType {
	case colliders.ColliderTypeStatic:
		c.staticPartitions.Remove(idx)
	default:
		c.partitions.Remove(idx)
	}

	delete(c.entities, entity)
	return true
}

func (c *collision) Update(entry *donburi.Entry) bool {
	entity := entry.Entity()

	idx, exists := c.entities[entity]
	if !exists {
		return false
	}

	collider := colliders.GetBoxCollider(entry)
	position := transform.GetPosition(entry)

	return c.partitions.Reinsert(idx, collider.Translate(position))
}

func (c *collision) Query(aabb geom.AABB) []donburi.Entity {
	c.queryCache = c.queryCache[:0]
	c.partitions.QueryAll(aabb, func(e donburi.Entity) bool {
		c.queryCache = append(c.queryCache, e)
		return true
	})
	return c.queryCache
}

func (c *collision) QueryStatic(aabb geom.AABB) []donburi.Entity {
	c.queryCache = c.queryCache[:0]
	c.staticPartitions.QueryAll(aabb, func(e donburi.Entity) bool {
		c.queryCache = append(c.queryCache, e)
		return true
	})
	return c.queryCache
}

func (c *collision) Simulate() {

}
