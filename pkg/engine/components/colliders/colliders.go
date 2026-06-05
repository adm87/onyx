package colliders

import (
	"github.com/adm87/onyx/pkg/engine/geom"
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
		IsStatic bool
		Layer    CollisionLayer
		Box      geom.AABB
	}
	Option func(*Options)
)

type QueryCallback func(*donburi.Entry, CollisionLayer, geom.AABB)

var (
	defaultCollisionInfo collisionInfo = collisionInfo{
		Enabled: true,
		Layer:   0,
	}
	defaultBox geom.AABB = geom.AABB{
		Min: geom.Vec2{X: 0, Y: 0},
		Max: geom.Vec2{X: 1, Y: 1},
	}
)

var (
	StaticType  = donburi.NewTag()
	DynamicType = donburi.NewTag()
)

var (
	Collision = donburi.NewComponentType[collisionInfo](defaultCollisionInfo)
	Box       = donburi.NewComponentType[geom.AABB](defaultBox)
)

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
		IsStatic: false,
		Layer:    0,
		Box:      defaultBox,
	}
}

func AsStatic() Option {
	return func(opts *Options) {
		opts.IsStatic = true
	}
}

func WithLayer(layer CollisionLayer) Option {
	return func(opts *Options) {
		opts.Layer = layer
	}
}

func WithBox(box geom.AABB) Option {
	return func(opts *Options) {
		opts.Box = box
	}
}

// QueryStatic iterates over all static colliders in the world, invoking the provided callback for each one.
func QueryStatic(world donburi.World, fn QueryCallback) {
	staticQuery.Each(world, func(entry *donburi.Entry) {
		layer := GetCollisionLayer(entry)
		box := GetBox(entry)
		fn(entry, layer, box)
	})
}

// QueryDynamic iterates over all dynamic colliders in the world, invoking the provided callback for each one.
func QueryDynamic(world donburi.World, fn QueryCallback) {
	dynamicQuery.Each(world, func(entry *donburi.Entry) {
		layer := GetCollisionLayer(entry)
		box := GetBox(entry)
		fn(entry, layer, box)
	})
}

// Query iterates over all colliders in the world, invoking the provided callback for each one.
func Query(world donburi.World, fn QueryCallback) {
	QueryStatic(world, fn)
	QueryDynamic(world, fn)
}

// QueryWith allows querying colliders using a custom query, invoking the provided callback for each matching entry.
func QueryWith(world donburi.World, q *donburi.Query, fn QueryCallback) {
	q.Each(world, func(entry *donburi.Entry) {
		layer := GetCollisionLayer(entry)
		box := GetBox(entry)
		fn(entry, layer, box)
	})
}

// NewBoxCollider creates a new entity with a box collider component, applying the provided options for configuration.
func NewBoxCollider(world donburi.World, options ...Option) *donburi.Entry {
	return AddBoxCollider(world.Entry(
		world.Create(
			Collision,
			Box,
		),
	), options...)
}

// AddBoxCollider adds a box collider component to an existing entity, applying the provided options for configuration.
func AddBoxCollider(entry *donburi.Entry, options ...Option) *donburi.Entry {
	opts := defaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	SetCollisionLayer(entry, opts.Layer)
	SetBox(entry, opts.Box)

	if opts.IsStatic {
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

// GetBox retrieves the box collider component from an entity, returning a default value if it does not exist.
func GetBox(entry *donburi.Entry) geom.AABB {
	if !entry.HasComponent(Box) {
		return defaultBox
	}
	return *Box.Get(entry)
}

// SetBox sets the box collider component for an entity, adding it if it does not already exist.
func SetBox(entry *donburi.Entry, box geom.AABB) {
	donburi.Add(entry, Box, &box)
}
