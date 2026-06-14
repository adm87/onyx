package transform

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type transformMatrix struct {
	dirty bool
	geom  ebiten.GeoM
}

type TransformModel struct {
	X, Y           float64
	ScaleX, ScaleY float64
	Rotation       float64
}

type (
	Options struct {
		Position geom.Vec2
		Scale    geom.Vec2
		Rotation float64
		Bounds   geom.AABB
	}
	Option func(*Options)
)

var (
	defaultPosition = geom.Vec2{X: 0, Y: 0}
	defaultScale    = geom.Vec2{X: 1, Y: 1}
	defaultRotation = 0.0
	defaultMatrix   = transformMatrix{dirty: true}
)

var (
	Transform = donburi.NewComponentType[TransformModel]()
	Bounds    = donburi.NewComponentType[geom.AABB](geom.AABB{})
	matrix    = donburi.NewComponentType[transformMatrix](defaultMatrix)
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

func WithBounds(bounds geom.AABB) Option {
	return func(opts *Options) {
		opts.Bounds = bounds
	}
}

// NewTransform creates a new entity with the necessary components for position, scale, rotation, and transformation matrix.
func NewTransform(ecs donburi.World, options ...Option) *donburi.Entry {
	return AddTransform(ecs.Entry(
		ecs.Create(
			Transform,
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

	SetTransform(entry, opts.Position, opts.Scale, opts.Rotation)
	SetLocalBounds(entry, &opts.Bounds)

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
	if !entry.HasComponent(Transform) {
		return defaultPosition
	}
	t := Transform.Get(entry)
	return geom.Vec2{X: t.X, Y: t.Y}
}

// SetPosition sets the position component for an entity, adding it if it does not already exist.
func SetPosition(entry *donburi.Entry, pos geom.Vec2) {
	if !entry.HasComponent(Transform) {
		return
	}

	t := Transform.Get(entry)
	if t.X == pos.X && t.Y == pos.Y {
		return
	}

	t.X, t.Y = pos.X, pos.Y
	markDirty(entry)
}

// GetScale retrieves the scale component from an entity, returning a default value if it does not exist.
func GetScale(entry *donburi.Entry) geom.Vec2 {
	if !entry.HasComponent(Transform) {
		return defaultScale
	}
	t := Transform.Get(entry)
	return geom.Vec2{X: t.ScaleX, Y: t.ScaleY}
}

// SetScale sets the scale component for an entity, adding it if it does not already exist.
func SetScale(entry *donburi.Entry, x, y float64) {
	if !entry.HasComponent(Transform) {
		return
	}

	t := Transform.Get(entry)
	if t.ScaleX == x && t.ScaleY == y {
		return
	}

	t.ScaleX, t.ScaleY = x, y
	markDirty(entry)
}

// GetRotation retrieves the rotation component from an entity, returning a default value if it does not exist.
func GetRotation(entry *donburi.Entry) float64 {
	if !entry.HasComponent(Transform) {
		return defaultRotation
	}
	t := Transform.Get(entry)
	return t.Rotation
}

// SetRotation sets the rotation component for an entity, adding it if it does not already exist.
func SetRotation(entry *donburi.Entry, rotation float64) {
	if !entry.HasComponent(Transform) {
		return
	}

	t := Transform.Get(entry)
	if t.Rotation == rotation {
		return
	}

	t.Rotation = rotation
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
		t := Transform.Get(entry)

		m.geom.Reset()
		m.geom.Scale(t.ScaleX, t.ScaleY)
		m.geom.Rotate(t.Rotation)
		m.geom.Translate(t.X, t.Y)

		m.dirty = false
	}
	return m.geom
}

func GetTransform(entry *donburi.Entry) *TransformModel {
	if !entry.HasComponent(Transform) {
		return &TransformModel{
			X:        defaultPosition.X,
			Y:        defaultPosition.Y,
			ScaleX:   defaultScale.X,
			ScaleY:   defaultScale.Y,
			Rotation: defaultRotation,
		}
	}
	return Transform.Get(entry)
}

func SetTransform(entry *donburi.Entry, position geom.Vec2, scale geom.Vec2, rotation float64) {
	donburi.Add(entry, Transform, &TransformModel{
		X:        position.X,
		Y:        position.Y,
		ScaleX:   scale.X,
		ScaleY:   scale.Y,
		Rotation: rotation,
	})
	markDirty(entry)
}

func GetLocalBounds(entry *donburi.Entry) *geom.AABB {
	if !entry.HasComponent(Bounds) {
		return &geom.AABB{}
	}
	return Bounds.Get(entry)
}

func SetLocalBounds(entry *donburi.Entry, bounds *geom.AABB) {
	donburi.Add(entry, Bounds, bounds)
}

func GetWorldBounds(entry *donburi.Entry) geom.AABB {
	if !entry.HasComponent(Bounds) {
		return geom.AABB{}
	}

	t := Transform.Get(entry)
	return Bounds.Get(entry).Translate(t.X, t.Y)
}
