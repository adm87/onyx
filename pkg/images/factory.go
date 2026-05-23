package images

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Options struct {
	Position geom.Vec2
	Rotation float64
	Scale    geom.Vec2
	Image    *ebiten.Image
	Anchor   geom.Vec2
	Color    color.RGBA
}

type Option func(*Options)

func defaultOptions() *Options {
	return &Options{
		Position: geom.Vec2{X: 0, Y: 0},
		Rotation: 0,
		Scale:    geom.Vec2{X: 1, Y: 1},
		Image:    nil,
		Anchor:   geom.Vec2{X: 0, Y: 0},
		Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}
}

func NewEntity(world donburi.World, opts ...Option) donburi.Entity {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	entity := world.Create(Components...)
	entry := world.Entry(entity)

	transform.SetPosition(entry, options.Position)
	transform.SetRotation(entry, options.Rotation)
	transform.SetScale(entry, options.Scale)

	if options.Image != nil {
		rendering.SetImage(entry, options.Image)
	}
	rendering.SetAnchor(entry, options.Anchor)
	rendering.SetColor(entry, options.Color)

	return entity
}

func WithPosition(x, y float64) Option {
	return func(opts *Options) {
		opts.Position = geom.Vec2{X: x, Y: y}
	}
}

func WithRotation(rotation float64) Option {
	return func(opts *Options) {
		opts.Rotation = rotation
	}
}

func WithScale(x, y float64) Option {
	return func(opts *Options) {
		opts.Scale = geom.Vec2{X: x, Y: y}
	}
}

func WithImage(img *ebiten.Image) Option {
	return func(opts *Options) {
		opts.Image = img
	}
}

func WithAnchor(x, y float64) Option {
	return func(opts *Options) {
		opts.Anchor = geom.Vec2{X: x, Y: y}
	}
}

func WithColor(color color.RGBA) Option {
	return func(opts *Options) {
		opts.Color = color
	}
}

func WithAlpha(alpha uint8) Option {
	return func(opts *Options) {
		opts.Color.A = alpha
	}
}
