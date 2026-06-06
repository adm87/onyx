package shapes

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type (
	Options struct {
		Position geom.Vec2
		Size     geom.Vec2
	}
	Option func(*Options)
)

var defaultBox geom.AABB = geom.AABB{
	Min: geom.Vec2{X: 0, Y: 0},
	Max: geom.Vec2{X: 1, Y: 1},
}

var AABB = donburi.NewComponentType[geom.AABB](defaultBox)

var queryAABB = donburi.NewQuery(
	filter.Contains(AABB),
)

func defaultOptions() Options {
	return Options{
		Position: defaultBox.Min,
		Size:     geom.Vec2{X: defaultBox.Max.X - defaultBox.Min.X, Y: defaultBox.Max.Y - defaultBox.Min.Y},
	}
}

func WithPosition(x, y float64) Option {
	return func(opts *Options) {
		opts.Position = geom.Vec2{X: x, Y: y}
	}
}

func WithSize(width, height float64) Option {
	return func(opts *Options) {
		opts.Size = geom.Vec2{X: width, Y: height}
	}
}

func WithBounds(min, max geom.Vec2) Option {
	return func(opts *Options) {
		opts.Position = min
		opts.Size = geom.Vec2{X: max.X - min.X, Y: max.Y - min.Y}
	}
}

func QueryAABB(ecs donburi.World, callback func(donburi.Entity)) {
	queryAABB.Each(ecs, func(entry *donburi.Entry) {
		callback(entry.Entity())
	})
}

func NewAABB(ecs donburi.World, options ...Option) *donburi.Entry {
	return AddAABB(ecs.Entry(
		ecs.Create(
			AABB,
		),
	), options...)
}

func AddAABB(entry *donburi.Entry, options ...Option) *donburi.Entry {
	opts := defaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	SetAABB(entry, geom.AABB{
		Min: opts.Position,
		Max: geom.Vec2{
			X: opts.Position.X + opts.Size.X,
			Y: opts.Position.Y + opts.Size.Y,
		},
	})

	return entry
}

func GetAABB(entry *donburi.Entry) geom.AABB {
	if !entry.HasComponent(AABB) {
		return defaultBox
	}
	return *AABB.Get(entry)
}

func SetAABB(entry *donburi.Entry, box geom.AABB) {
	donburi.Add(entry, AABB, &box)
}

func IsAABB(entry *donburi.Entry) bool {
	return entry.HasComponent(AABB)
}
