package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type CollisionEvent struct {
	EntityA donburi.Entity
	EntityB donburi.Entity
}

type Collision interface {
}

type (
	collisionMask    map[colliders.CollisionLayer]colliders.CollisionLayer
	collisionPairing map[collisionPair]struct{}
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
}

func newCollision() *collision {
	return &collision{}
}
