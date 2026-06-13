package transform

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type transformMatrix struct {
	dirty bool
	geom  ebiten.GeoM
}

type (
	Options struct {
		Position geom.Vec2
		Scale    geom.Vec2
		Rotation float64
	}
	Option func(*Options)
)

type QueryCallback func(*donburi.Entry, geom.Vec2, geom.Vec2, float64)

var (
	defaultPosition = geom.Vec2{X: 0, Y: 0}
	defaultScale    = geom.Vec2{X: 1, Y: 1}
	defaultRotation = 0.0
	defaultMatrix   = transformMatrix{dirty: true}
)

var (
	Position = donburi.NewComponentType[geom.Vec2](defaultPosition)
	Scale    = donburi.NewComponentType[geom.Vec2](defaultScale)
	Rotation = donburi.NewComponentType[float64](defaultRotation)
	matrix   = donburi.NewComponentType[transformMatrix](defaultMatrix)
)

var query = donburi.NewQuery(
	filter.Contains(
		Position,
		Scale,
		Rotation,
		matrix,
	),
)

func defaultOptions() Options {
	return Options{
		Position: defaultPosition,
		Scale:    defaultScale,
		Rotation: defaultRotation,
	}
}

func WithPosition(x, y float64) Option {
	return func(opts *Options) {
		opts.Position = geom.Vec2{X: x, Y: y}
	}
}

func WithScale(x, y float64) Option {
	return func(opts *Options) {
		opts.Scale = geom.Vec2{X: x, Y: y}
	}
}

func WithRotation(rotation float64) Option {
	return func(opts *Options) {
		opts.Rotation = rotation
	}
}

// Query iterates over all entities with the necessary components for transformation and applies the provided function to each entry.
func Query(ecs donburi.World, fn QueryCallback) {
	query.Each(ecs, func(entry *donburi.Entry) {
		position := GetPosition(entry)
		scale := GetScale(entry)
		rotation := GetRotation(entry)
		fn(entry, position, scale, rotation)
	})
}

// QueryWith allows you to specify a custom query to iterate over entities with the necessary components for transformation
// and applies the provided function to each entry.
func QueryWith(ecs donburi.World, q *donburi.Query, fn QueryCallback) {
	q.Each(ecs, func(entry *donburi.Entry) {
		position := GetPosition(entry)
		scale := GetScale(entry)
		rotation := GetRotation(entry)
		fn(entry, position, scale, rotation)
	})
}

// NewTransform creates a new entity with the necessary components for position, scale, rotation, and transformation matrix.
func NewTransform(ecs donburi.World, options ...Option) *donburi.Entry {
	return AddTransform(ecs.Entry(
		ecs.Create(
			Position,
			Scale,
			Rotation,
			matrix,
		),
	), options...)
}

// AddTransform adds the necessary components for position, scale, rotation, and transformation matrix to an existing entity.
func AddTransform(entry *donburi.Entry, options ...Option) *donburi.Entry {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}

	SetPosition(entry, opts.Position)
	SetScale(entry, opts.Scale.X, opts.Scale.Y)
	SetRotation(entry, opts.Rotation)

	m := defaultMatrix
	donburi.Add(entry, matrix, &m)

	return entry
}

func markDirty(entry *donburi.Entry) {
	if entry.HasComponent(matrix) {
		m := matrix.Get(entry)
		m.dirty = true
	}
}

// GetPosition retrieves the position component from an entity, returning a default value if it does not exist.
func GetPosition(entry *donburi.Entry) geom.Vec2 {
	if !entry.HasComponent(Position) {
		return defaultPosition
	}
	return *Position.Get(entry)
}

// SetPosition sets the position component for an entity, adding it if it does not already exist.
func SetPosition(entry *donburi.Entry, pos geom.Vec2) {
	p := GetPosition(entry)
	if p == pos {
		return
	}
	donburi.Add(entry, Position, &pos)
	markDirty(entry)
}

// GetScale retrieves the scale component from an entity, returning a default value if it does not exist.
func GetScale(entry *donburi.Entry) geom.Vec2 {
	if !entry.HasComponent(Scale) {
		return defaultScale
	}
	return *Scale.Get(entry)
}

// SetScale sets the scale component for an entity, adding it if it does not already exist.
func SetScale(entry *donburi.Entry, x, y float64) {
	s := GetScale(entry)
	if s.X == x && s.Y == y {
		return
	}
	scale := geom.Vec2{X: x, Y: y}
	donburi.Add(entry, Scale, &scale)
	markDirty(entry)
}

// GetRotation retrieves the rotation component from an entity, returning a default value if it does not exist.
func GetRotation(entry *donburi.Entry) float64 {
	if !entry.HasComponent(Rotation) {
		return defaultRotation
	}
	return *Rotation.Get(entry)
}

// SetRotation sets the rotation component for an entity, adding it if it does not already exist.
func SetRotation(entry *donburi.Entry, rotation float64) {
	r := GetRotation(entry)
	if r == rotation {
		return
	}
	donburi.Add(entry, Rotation, &rotation)
	markDirty(entry)
}

// GetMatrix retrieves the transformation matrix for an entity,
// calculating it if necessary based on the position, scale, and rotation components.
func GetMatrix(entry *donburi.Entry) ebiten.GeoM {
	if !entry.HasComponent(matrix) {
		return ebiten.GeoM{}
	}

	m := matrix.Get(entry)
	if m.dirty {
		position := GetPosition(entry)
		scale := GetScale(entry)
		rotation := GetRotation(entry)

		m.geom.Reset()
		m.geom.Scale(scale.X, scale.Y)
		m.geom.Rotate(rotation)
		m.geom.Translate(position.X, position.Y)

		m.dirty = false
	}
	return m.geom
}
