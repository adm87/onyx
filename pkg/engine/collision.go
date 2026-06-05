package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/spatialhash"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type CollisionEvent struct {
	EntityA donburi.Entity
	EntityB donburi.Entity
}

// Collision defines the interface for managing collision detection and events within the game engine.
//
// This system is a glorified overlap detector, and does not resolve collisions between entities.
// Instead, it tracks the state of any overlap between entities and provides events for when overlaps begin, end, or persist.
//
// Static and dynamic colliders are indexed separately. Only dynamic colliders will initiate collision checks,
// but will check against both static and dynamic colliders. This means static colliders can overlap without triggering events.
type Collision interface {
	FlagLayers(a, b colliders.CollisionLayer)       // Flags two collision layers to allow interactions between them during collision checks.
	UnflagLayers(a, b colliders.CollisionLayer)     // Unflags two collision layers to prevent interactions between them during collision checks.
	CheckLayers(a, b colliders.CollisionLayer) bool // Checks if two collision layers are flagged to interact with each other.

	AddCollisionEnter(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent))
	AddCollisionExit(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent))
	AddCollisionStay(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent))
	RemoveCollisionEnter(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent))
	RemoveCollisionExit(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent))
	RemoveCollisionStay(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent))

	QueryAll(aabb geom.AABB, callback func(entity donburi.Entity))
	QueryStatic(aabb geom.AABB, callback func(entity donburi.Entity))
	QueryDynamic(aabb geom.AABB, callback func(entity donburi.Entity))
}

type (
	collisionIndexing map[donburi.Entity]spatialhash.SpatialIndex
	collisionMask     map[colliders.CollisionLayer]colliders.CollisionLayer
	collisionPairing  map[collisionPair]struct{}
)

type collisionEvents struct {
	enter *events.EventType[CollisionEvent]
	exit  *events.EventType[CollisionEvent]
	stay  *events.EventType[CollisionEvent]
}

type collisionPair [2]donburi.Entity

func newCollisionPair(a, b donburi.Entity) collisionPair {
	if a < b {
		return collisionPair{a, b}
	}
	return collisionPair{b, a}
}

func (p collisionPair) Low() donburi.Entity {
	return p[0]
}

func (p collisionPair) High() donburi.Entity {
	return p[1]
}

type collision struct {
	static  *spatialhash.SpatialHash[donburi.Entity]
	dynamic *spatialhash.SpatialHash[donburi.Entity]
	events  *collisionEvents

	indexing collisionIndexing
	masks    collisionMask

	currentPairs  collisionPairing
	previousPairs collisionPairing
}

func newCollision() *collision {
	return &collision{
		static: spatialhash.New[donburi.Entity](
			spatialhash.WithResolutions(16),
			spatialhash.WithCapacity(100),
		),
		dynamic: spatialhash.New[donburi.Entity](
			spatialhash.WithResolutions(16),
			spatialhash.WithCapacity(100),
		),
		events: &collisionEvents{
			enter: events.NewEventType[CollisionEvent](),
			exit:  events.NewEventType[CollisionEvent](),
			stay:  events.NewEventType[CollisionEvent](),
		},
		indexing:      make(collisionIndexing),
		masks:         make(collisionMask),
		currentPairs:  make(collisionPairing, 100),
		previousPairs: make(collisionPairing, 100),
	}
}

func (c *collision) add(entry *donburi.Entry, aabb geom.AABB) bool {
	entity := entry.Entity()

	if _, exists := c.indexing[entity]; exists {
		return false // Entity is already indexed, cannot add again
	}

	var index spatialhash.SpatialIndex
	var ok bool

	if colliders.IsStatic(entry) {
		index, ok = c.static.Insert(entity, aabb)
	} else {
		index, ok = c.dynamic.Insert(entity, aabb)
	}

	if ok {
		c.indexing[entity] = index
	}

	return ok
}

func (c *collision) remove(entry *donburi.Entry) bool {
	entity := entry.Entity()

	index, exists := c.indexing[entity]
	if !exists {
		return true // Entity is not indexed, consider it removed
	}

	delete(c.indexing, entity)

	if colliders.IsStatic(entry) {
		c.static.Remove(index)
	} else {
		c.dynamic.Remove(index)
	}

	return true
}

func (c *collision) update(entry *donburi.Entry, aabb geom.AABB) bool {
	entity := entry.Entity()

	index, exists := c.indexing[entity]
	if !exists {
		return false // Entity is not indexed, cannot update
	}

	if colliders.IsStatic(entry) {
		return c.static.Reinsert(index, aabb)
	}

	return c.dynamic.Reinsert(index, aabb)
}

func (c *collision) FlagLayers(a, b colliders.CollisionLayer) {
	c.masks[a] |= b
	c.masks[b] |= a
}

func (c *collision) UnflagLayers(a, b colliders.CollisionLayer) {
	c.masks[a] &^= b
	c.masks[b] &^= a
}

func (c *collision) CheckLayers(a, b colliders.CollisionLayer) bool {
	if a == b {
		return true
	}
	return (c.masks[a] & b) != 0
}

func (c *collision) AddCollisionEnter(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent)) {
	c.events.enter.Subscribe(ecs, callback)
}

func (c *collision) AddCollisionExit(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent)) {
	c.events.exit.Subscribe(ecs, callback)
}

func (c *collision) AddCollisionStay(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent)) {
	c.events.stay.Subscribe(ecs, callback)
}

func (c *collision) RemoveCollisionEnter(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent)) {
	c.events.enter.Unsubscribe(ecs, callback)
}

func (c *collision) RemoveCollisionExit(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent)) {
	c.events.exit.Unsubscribe(ecs, callback)
}

func (c *collision) RemoveCollisionStay(ecs donburi.World, callback func(ecs donburi.World, event CollisionEvent)) {
	c.events.stay.Unsubscribe(ecs, callback)
}

func (c *collision) QueryAll(aabb geom.AABB, callback func(entity donburi.Entity)) {
	c.QueryStatic(aabb, callback)
	c.QueryDynamic(aabb, callback)
}

func (c *collision) QueryStatic(aabb geom.AABB, callback func(entity donburi.Entity)) {
	c.static.QueryAll(aabb, func(entity donburi.Entity) bool {
		callback(entity)
		return true
	})
}

func (c *collision) QueryDynamic(aabb geom.AABB, callback func(entity donburi.Entity)) {
	c.dynamic.QueryAll(aabb, func(entity donburi.Entity) bool {
		callback(entity)
		return true
	})
}

func (c *collision) checkCollision(ecs donburi.World) error {
	c.currentPairs, c.previousPairs = c.previousPairs, c.currentPairs
	clear(c.currentPairs)

	// Broad phase: Collect potential collision pairs based on spatial hashing
	colliders.QueryEnabledDynamic(ecs, func(entry *donburi.Entry, cl colliders.CollisionLayer, aabb geom.AABB) {
		aabb = aabb.Translate(transform.GetPosition(entry))
		entity := entry.Entity()

		c.static.QueryAll(aabb, func(otherEntity donburi.Entity) bool {
			return c.validatePair(ecs, entity, otherEntity, aabb, cl)
		})
		c.dynamic.QueryAll(aabb, func(otherEntity donburi.Entity) bool {
			return c.validatePair(ecs, entity, otherEntity, aabb, cl)
		})
	})

	var event CollisionEvent

	// Narrow phase: Determine collision events based on current and previous pairs
	for pair := range c.currentPairs {
		event = CollisionEvent{EntityA: pair.Low(), EntityB: pair.High()}
		if _, existed := c.previousPairs[pair]; existed {
			c.events.stay.Publish(ecs, event)
		} else {
			c.events.enter.Publish(ecs, event)
		}
	}
	for pair := range c.previousPairs {
		if _, exists := c.currentPairs[pair]; !exists {
			event = CollisionEvent{EntityA: pair.Low(), EntityB: pair.High()}
			c.events.exit.Publish(ecs, event)
		}
	}

	c.events.enter.ProcessEvents(ecs)
	c.events.exit.ProcessEvents(ecs)
	c.events.stay.ProcessEvents(ecs)

	return nil
}

func (c *collision) validatePair(ecs donburi.World, entityA, entityB donburi.Entity, boxA geom.AABB, layerA colliders.CollisionLayer) bool {
	if entityA == entityB {
		return true // Skip self-collision
	}
	otherEntry := ecs.Entry(entityB)

	if !colliders.IsCollisionEnabled(otherEntry) {
		return true // Collision is not enabled for this entity, skip
	}

	otherCL := colliders.GetCollisionLayer(otherEntry)
	otherAABB := colliders.GetAABB(otherEntry).Translate(transform.GetPosition(otherEntry))

	if !c.CheckLayers(layerA, otherCL) {
		return true // Layers are not flagged to interact, skip
	}

	if !boxA.Intersects(otherAABB) {
		return true // AABBs do not intersect, skip
	}

	pair := newCollisionPair(entityA, entityB)
	c.currentPairs[pair] = struct{}{}

	return true
}
