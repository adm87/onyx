package images

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type Options struct {
	position geom.Vec2
	anchor   geom.Vec2
	scale    geom.Vec2
	layer    int
	zIndex   int
	visible  bool
	color    color.RGBA
	ref      *ebiten.Image
}

type Option func(*Options)

var query = donburi.NewQuery(
	filter.Contains(rendering.Image),
)

func defaultOptions() Options {
	return Options{
		position: geom.Vec2{X: 0, Y: 0},
		anchor:   geom.Vec2{X: 0, Y: 0},
		layer:    0,
		zIndex:   0,
		visible:  true,
		color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
		scale:    geom.Vec2{X: 1, Y: 1},
	}
}

func CreateImage(world donburi.World, opts ...Option) *donburi.Entry {
	entity := world.Create(
		rendering.Image,
		rendering.Renderer,
		rendering.Anchor,
		rendering.Color,
		transform.Position,
		transform.Scale,
	)
	entry := world.Entry(entity)

	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	rendering.SetImage(entry, options.ref)
	rendering.SetAnchor(entry, options.anchor)
	rendering.SetLayer(entry, options.layer)
	rendering.SetZIndex(entry, options.zIndex)
	rendering.SetVisible(entry, options.visible)
	rendering.SetColor(entry, options.color)

	transform.SetPosition(entry, options.position)
	transform.SetScale(entry, options.scale)

	return entry
}

func ForEach(world donburi.World, fn func(*donburi.Entry)) {
	query.Each(world, func(entry *donburi.Entry) {
		fn(entry)
	})
}

func WithRef(img *ebiten.Image) Option {
	return func(opts *Options) {
		opts.ref = img
	}
}

func WithAnchor(x, y float64) Option {
	return func(opts *Options) {
		opts.anchor = geom.Vec2{X: x, Y: y}
	}
}

func WithLayer(layer int) Option {
	return func(opts *Options) {
		opts.layer = layer
	}
}

func WithZIndex(zIndex int) Option {
	return func(opts *Options) {
		opts.zIndex = zIndex
	}
}

func WithColor(col color.RGBA) Option {
	return func(opts *Options) {
		opts.color = col
	}
}

func WithPosition(x, y float64) Option {
	return func(opts *Options) {
		opts.position = geom.Vec2{X: x, Y: y}
	}
}

func WithScale(x, y float64) Option {
	return func(opts *Options) {
		opts.scale = geom.Vec2{X: x, Y: y}
	}
}

func WithVisible(visible bool) Option {
	return func(opts *Options) {
		opts.visible = visible
	}
}
