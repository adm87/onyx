package collision

import (
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
)

type CollisionType uint8

type CollisionModel struct {
	Type    CollisionType
	Enabled bool
}

type CollisionOptions struct {
	Type    CollisionType
	Enabled bool
	Bounds  *geom.AABB
}

type CollisionOption func(*CollisionOptions)

const (
	CollisionTypeStatic CollisionType = iota
	CollisionTypeDynamic
)

var (
	Collision      = donburi.NewComponentType[CollisionModel]()
	CollisionIndex = donburi.NewComponentType[uint64]()
	Collider       = donburi.NewComponentType[geom.AABB]()
)

func WithCollisionType(collisionType CollisionType) CollisionOption {
	return func(o *CollisionOptions) {
		o.Type = collisionType
	}
}

func WithCollisionEnabled(enabled bool) CollisionOption {
	return func(o *CollisionOptions) {
		o.Enabled = enabled
	}
}

func WithCollisionBounds(min, max geom.Vec2) CollisionOption {
	return func(o *CollisionOptions) {
		o.Bounds = &geom.AABB{Min: min, Max: max}
	}
}

func AddCollision(entry *donburi.Entry, opts ...CollisionOption) {
	options := &CollisionOptions{
		Type:    CollisionTypeStatic,
		Enabled: true,
	}
	for _, opt := range opts {
		opt(options)
	}

	donburi.Add(entry, Collision, &CollisionModel{
		Type:    options.Type,
		Enabled: options.Enabled,
	})

	var bounds geom.AABB

	if options.Bounds != nil {
		bounds = *options.Bounds
	} else {
		bounds = transform.GetBounds(entry)
	}

	donburi.Add(entry, Collider, &bounds)
}

func HasCollision(entry *donburi.Entry) bool {
	return entry.HasComponent(Collision)
}

func GetCollisionType(entry *donburi.Entry) CollisionType {
	if !entry.HasComponent(Collision) {
		return CollisionTypeStatic
	}
	return Collision.Get(entry).Type
}

func IsEnabled(entry *donburi.Entry) bool {
	if !entry.HasComponent(Collision) {
		return true
	}
	return Collision.Get(entry).Enabled
}

func SetEnabled(entry *donburi.Entry, enabled bool) {
	if !entry.HasComponent(Collision) {
		return
	}
	collision := Collision.Get(entry)
	collision.Enabled = enabled
}

func GetCollider(entry *donburi.Entry) geom.AABB {
	if !entry.HasComponent(Collider) {
		return geom.AABB{}
	}
	return *Collider.Get(entry)
}

func GetWorldCollider(entry *donburi.Entry) geom.AABB {
	bounds := GetCollider(entry)
	matrix := transform.GetMatrix(entry)

	x1, y1 := matrix.Apply(bounds.Min.X, bounds.Min.Y)
	x2, y2 := matrix.Apply(bounds.Max.X, bounds.Max.Y)

	return geom.AABB{
		Min: geom.Vec2{
			X: min(x1, x2),
			Y: min(y1, y2),
		},
		Max: geom.Vec2{
			X: max(x1, x2),
			Y: max(y1, y2),
		},
	}
}

func SetBoxCollider(entry *donburi.Entry, aabb geom.AABB) {
	donburi.Add(entry, Collider, &aabb)
}

func GetCollisionIndex(entry *donburi.Entry) uint64 {
	if !entry.HasComponent(CollisionIndex) {
		return 0
	}
	return *CollisionIndex.Get(entry)
}

func SetCollisionIndex(entry *donburi.Entry, index uint64) {
	donburi.Add(entry, CollisionIndex, &index)
}
