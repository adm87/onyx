package collision

import (
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
)

type CollisionOptions struct {
	BoxCollider *geom.AABB
}

type Option func(*CollisionOptions)

type CollisionInfo struct {
	Enabled bool
}

var (
	Collision   = donburi.NewComponentType[CollisionInfo]()
	BoxCollider = donburi.NewComponentType[geom.AABB]()
)

func WithBoxCollider(aabb geom.AABB) Option {
	return func(opts *CollisionOptions) {
		opts.BoxCollider = &aabb
	}
}

func AddCollision(entry *donburi.Entry, options ...Option) {
	opts := CollisionOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	SetCollision(entry, &CollisionInfo{Enabled: true})

	if opts.BoxCollider != nil {
		SetBoxCollider(entry, opts.BoxCollider)
	}
}

func GetCollision(entry *donburi.Entry) *CollisionInfo {
	if !entry.HasComponent(Collision) {
		return nil
	}
	return Collision.Get(entry)
}

func SetCollision(entry *donburi.Entry, info *CollisionInfo) {
	donburi.Add(entry, Collision, info)
}

func GetBoxCollider(entry *donburi.Entry) geom.AABB {
	if !entry.HasComponent(BoxCollider) {
		return *transform.GetLocalBounds(entry)
	}
	return *BoxCollider.Get(entry)
}

func SetBoxCollider(entry *donburi.Entry, aabb *geom.AABB) {
	donburi.Add(entry, BoxCollider, aabb)
}
