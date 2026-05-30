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

	Add(entry *donburi.Entry) bool
	Update(entry *donburi.Entry) bool

	Query(aabb geom.AABB) []donburi.Entity
}

type collision struct {
	partitioning *partitioning.SpatialHash[donburi.Entity]
	entities     map[donburi.Entity]partitioning.SpatialIndex

	queryCache []donburi.Entity
}

func newCollision() *collision {
	return &collision{
		partitioning: partitioning.NewSpatialHash[donburi.Entity](64, 32, 64, 128),
		entities:     make(map[donburi.Entity]partitioning.SpatialIndex),
	}
}

func (c *collision) Partitioning() *partitioning.SpatialHash[donburi.Entity] {
	return c.partitioning
}

func (c *collision) Add(entry *donburi.Entry) bool {
	entity := entry.Entity()

	if _, exists := c.entities[entity]; exists {
		return false
	}

	collider := colliders.GetBoxCollider(entry)
	position := transform.GetPosition(entry)

	idx, ok := c.partitioning.Insert(entity, collider.Translate(position))
	if !ok {
		return false
	}

	c.entities[entity] = idx
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

	c.partitioning.Remove(idx)

	newIdx, ok := c.partitioning.Insert(entity, collider.Translate(position))
	if !ok {
		delete(c.entities, entity)
		return false
	}

	c.entities[entity] = newIdx
	return true
}

func (c *collision) Query(aabb geom.AABB) []donburi.Entity {
	c.queryCache = c.queryCache[:0]
	c.partitioning.QueryAll(aabb, func(e donburi.Entity) bool {
		c.queryCache = append(c.queryCache, e)
		return true
	})
	return c.queryCache
}
