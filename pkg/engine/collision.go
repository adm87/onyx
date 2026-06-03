package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/spatialhash"
	"github.com/yohamta/donburi"
)

type collisionPair struct {
	lo donburi.Entity
	hi donburi.Entity
}

func newCollisionPair(a, b donburi.Entity) collisionPair {
	return collisionPair{lo: min(a, b), hi: max(a, b)}
}

type CollisionEvent struct {
	A donburi.Entity
	B donburi.Entity
}

type Collision interface {
	Partitioning() *spatialhash.SpatialHash[donburi.Entity]
	StaticPartitioning() *spatialhash.SpatialHash[donburi.Entity]

	Add(entry *donburi.Entry) bool
	Remove(entry *donburi.Entry) bool
	Update(entry *donburi.Entry) bool
	Simulate(world donburi.World)

	Query(aabb geom.AABB) []donburi.Entity
	QueryStatic(aabb geom.AABB) []donburi.Entity
}

type collision struct {
	staticPartitions *spatialhash.SpatialHash[donburi.Entity]
	partitions       *spatialhash.SpatialHash[donburi.Entity]

	queryCache []donburi.Entity
	currPairs  []collisionPair

	prevPairs map[collisionPair]struct{}
	entities  map[donburi.Entity]spatialhash.SpatialIndex

	logger Logger
}

func newCollision(logger Logger) *collision {
	return &collision{
		staticPartitions: spatialhash.New[donburi.Entity](
			spatialhash.WithCapacity(100),
			spatialhash.WithResolutions(8),
			spatialhash.WithPadding(1, 1, 0, 0),
		),
		partitions: spatialhash.New[donburi.Entity](
			spatialhash.WithCapacity(100),
			spatialhash.WithResolutions(8),
		),
		queryCache: make([]donburi.Entity, 0, 100),
		currPairs:  make([]collisionPair, 0, 100),
		prevPairs:  make(map[collisionPair]struct{}, 100),
		entities:   make(map[donburi.Entity]spatialhash.SpatialIndex),
		logger:     logger,
	}
}

func (c *collision) StaticPartitioning() *spatialhash.SpatialHash[donburi.Entity] {
	return c.staticPartitions
}

func (c *collision) Partitioning() *spatialhash.SpatialHash[donburi.Entity] {
	return c.partitions
}

func (c *collision) Add(entry *donburi.Entry) bool {
	entity := entry.Entity()

	if _, exists := c.entities[entity]; exists {
		return false
	}

	collider := colliders.GetBoxCollider(entry)
	position := transform.GetPosition(entry)

	var idx spatialhash.SpatialIndex
	var ok bool

	switch {
	case colliders.IsStatic(entry):
		idx, ok = c.staticPartitions.Insert(entity, collider.Translate(position))
	case colliders.IsDynamic(entry), colliders.IsKinematic(entry):
		idx, ok = c.partitions.Insert(entity, collider.Translate(position))
	default:
		c.logger.Warn("cannot add collider: entity has collider component but no type tag, entity: %d", entity)
		return false
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

	switch {
	case colliders.IsStatic(entry):
		c.staticPartitions.Remove(idx)
	case colliders.IsDynamic(entry), colliders.IsKinematic(entry):
		c.partitions.Remove(idx)
	default:
		c.logger.Warn("cannot remove collider: entity has collider component but no type tag, entity: %d", entity)
		return false
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

func (c *collision) Simulate(world donburi.World) {
	c.currPairs = c.currPairs[:0]

	// Broad phase: find potential collision pairs
	colliders.DynamicColliderQuery.Each(world, func(entry *donburi.Entry) {
		entity := entry.Entity()

		position := transform.GetPosition(entry)
		collider := colliders.GetBoxCollider(entry).Translate(position)

		c.staticPartitions.QueryAll(collider, func(e donburi.Entity) bool {
			c.currPairs = append(c.currPairs, newCollisionPair(entity, e))
			return true
		})
		c.partitions.QueryAll(collider, func(e donburi.Entity) bool {
			if e <= entity {
				return true
			}
			c.currPairs = append(c.currPairs, newCollisionPair(entity, e))
			return true
		})
	})
}
