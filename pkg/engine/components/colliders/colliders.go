package colliders

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// CollisionLayer represents a layer for collision detection,
// allowing entities to be categorized into different groups for collision purposes.
type CollisionLayer uint32

type collisionInfo struct {
	Enabled bool
	Layer   CollisionLayer
}

type (
	Options struct {
		AsStatic bool
		Layer    CollisionLayer
	}
	Option func(*Options)
)

type QueryCallback func(*donburi.Entry, CollisionLayer)

var defaultCollisionInfo collisionInfo = collisionInfo{
	Enabled: true,
	Layer:   0,
}

var (
	StaticType  = donburi.NewTag()
	DynamicType = donburi.NewTag()
)

var Collision = donburi.NewComponentType[collisionInfo](defaultCollisionInfo)

var (
	staticQuery = donburi.NewQuery(
		filter.Contains(
			Collision,
			StaticType,
		),
	)
	dynamicQuery = donburi.NewQuery(
		filter.Contains(
			Collision,
			DynamicType,
		),
	)
)

func defaultOptions() Options {
	return Options{
		AsStatic: false,
		Layer:    0,
	}
}

func AsStatic() Option {
	return func(opts *Options) {
		opts.AsStatic = true
	}
}

func WithLayer(layer CollisionLayer) Option {
	return func(opts *Options) {
		opts.Layer = layer
	}
}

// Query iterates over all colliders in the world, invoking the provided callback for each one.
func Query(ecs donburi.World, fn QueryCallback) {
	QueryStatic(ecs, fn)
	QueryDynamic(ecs, fn)
}

// QueryWith allows querying colliders using a custom query, invoking the provided callback for each matching entry.
func QueryWith(ecs donburi.World, q *donburi.Query, fn QueryCallback) {
	q.Each(ecs, func(entry *donburi.Entry) {
		layer := GetCollisionLayer(entry)
		fn(entry, layer)
	})
}

// QueryEnabled iterates over all colliders in the world that have collision enabled, invoking the provided callback for each one.
func QueryEnabled(ecs donburi.World, fn QueryCallback) {
	Query(ecs, func(entry *donburi.Entry, layer CollisionLayer) {
		if IsCollisionEnabled(entry) {
			fn(entry, layer)
		}
	})
}

// QueryEnabledWith allows querying colliders using a custom query and invokes the provided callback for each matching entry that has collision enabled.
func QueryEnabledWith(ecs donburi.World, q *donburi.Query, fn QueryCallback) {
	QueryWith(ecs, q, func(entry *donburi.Entry, layer CollisionLayer) {
		if IsCollisionEnabled(entry) {
			fn(entry, layer)
		}
	})
}

// QueryStatic iterates over all static colliders in the world, invoking the provided callback for each one.
func QueryStatic(ecs donburi.World, fn QueryCallback) {
	staticQuery.Each(ecs, func(entry *donburi.Entry) {
		layer := GetCollisionLayer(entry)
		fn(entry, layer)
	})
}

// QueryEnabledStatic iterates over all static colliders in the world that have collision enabled, invoking the provided callback for each one.
func QueryEnabledStatic(ecs donburi.World, fn QueryCallback) {
	QueryStatic(ecs, func(entry *donburi.Entry, layer CollisionLayer) {
		if IsCollisionEnabled(entry) {
			fn(entry, layer)
		}
	})
}

// QueryDynamic iterates over all dynamic colliders in the world, invoking the provided callback for each one.
func QueryDynamic(ecs donburi.World, fn QueryCallback) {
	dynamicQuery.Each(ecs, func(entry *donburi.Entry) {
		layer := GetCollisionLayer(entry)
		fn(entry, layer)
	})
}

// QueryEnabledDynamic iterates over all dynamic colliders in the world that have collision enabled, invoking the provided callback for each one.
func QueryEnabledDynamic(ecs donburi.World, fn QueryCallback) {
	QueryDynamic(ecs, func(entry *donburi.Entry, layer CollisionLayer) {
		if IsCollisionEnabled(entry) {
			fn(entry, layer)
		}
	})
}

// NewCollider creates a new entity with a collider component, applying the provided options for configuration.
func NewCollider(ecs donburi.World, options ...Option) *donburi.Entry {
	return AddCollider(ecs.Entry(
		ecs.Create(
			Collision,
		),
	), options...)
}

// AddCollider adds a collider component to an existing entity, applying the provided options for configuration.
func AddCollider(entry *donburi.Entry, options ...Option) *donburi.Entry {
	opts := defaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	SetCollisionLayer(entry, opts.Layer)

	if opts.AsStatic {
		entry.AddComponent(StaticType)
	} else {
		entry.AddComponent(DynamicType)
	}

	return entry
}

// IsStatic checks if the entity has a static collider component, indicating that it is a static collider.
func IsStatic(entry *donburi.Entry) bool {
	return entry.HasComponent(StaticType)
}

// IsDynamic checks if the entity has a dynamic collider component, indicating that it is a dynamic collider.
func IsDynamic(entry *donburi.Entry) bool {
	return entry.HasComponent(DynamicType)
}

// GetCollisionLayer retrieves the collision layer of an entity,
// returning a default value if the entity does not have a collision component.
func GetCollisionLayer(entry *donburi.Entry) CollisionLayer {
	if !entry.HasComponent(Collision) {
		return defaultCollisionInfo.Layer
	}
	return Collision.Get(entry).Layer
}

// SetCollisionLayer sets the collision layer for an entity, adding a collision component if it does not already exist.
func SetCollisionLayer(entry *donburi.Entry, layer CollisionLayer) {
	if !entry.HasComponent(Collision) {
		c := defaultCollisionInfo
		c.Layer = layer
		donburi.Add(entry, Collision, &c)
		return
	}
	collision := Collision.Get(entry)
	collision.Layer = layer
}

// RemoveCollider removes the collider component from an entity, effectively removing it from collision detection.
func IsCollisionEnabled(entry *donburi.Entry) bool {
	if !entry.HasComponent(Collision) {
		return defaultCollisionInfo.Enabled
	}
	return Collision.Get(entry).Enabled
}

// SetCollisionEnabled enables or disables collision for an entity by setting the Enabled field in the collision component.
func SetCollisionEnabled(entry *donburi.Entry, enabled bool) {
	if !entry.HasComponent(Collision) {
		c := defaultCollisionInfo
		c.Enabled = enabled
		donburi.Add(entry, Collision, &c)
		return
	}
	collision := Collision.Get(entry)
	collision.Enabled = enabled
}
