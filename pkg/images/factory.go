package images

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images/components"
	"github.com/yohamta/donburi"
)

type ImageOptions struct {
	position geom.Vec2
	anchor   geom.Vec2
	scale    geom.Vec2
	layer    int
	zIndex   int
	visible  bool
	color    color.RGBA
	ref      engine.FilePath
}

type ImageOption func(*ImageOptions)

func defaultImageOptions() ImageOptions {
	return ImageOptions{
		position: geom.Vec2{X: 0, Y: 0},
		anchor:   geom.Vec2{X: 0, Y: 0},
		layer:    0,
		zIndex:   0,
		visible:  true,
		color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
		scale:    geom.Vec2{X: 1, Y: 1},
	}
}

func CreateImageEntity(world donburi.World, opts ...ImageOption) *donburi.Entry {
	entity := world.Create(
		components.Image,
		rendering.Renderer,
		rendering.Anchor,
		rendering.Color,
		transform.Position,
		transform.Scale,
	)
	entry := world.Entry(entity)

	options := defaultImageOptions()
	for _, opt := range opts {
		opt(&options)
	}

	components.SetImageRef(entry, options.ref)

	rendering.SetAnchor(entry, options.anchor)
	rendering.SetLayer(entry, options.layer)
	rendering.SetZIndex(entry, options.zIndex)
	rendering.SetVisible(entry, options.visible)
	rendering.SetColor(entry, options.color)

	transform.SetPosition(entry, options.position)
	transform.SetScale(entry, options.scale)

	return entry
}

func WithRef(ref engine.FilePath) ImageOption {
	return func(opts *ImageOptions) {
		opts.ref = ref
	}
}

func WithAnchor(x, y float64) ImageOption {
	return func(opts *ImageOptions) {
		opts.anchor = geom.Vec2{X: x, Y: y}
	}
}

func WithLayer(layer int) ImageOption {
	return func(opts *ImageOptions) {
		opts.layer = layer
	}
}

func WithZIndex(zIndex int) ImageOption {
	return func(opts *ImageOptions) {
		opts.zIndex = zIndex
	}
}

func WithColor(col color.RGBA) ImageOption {
	return func(opts *ImageOptions) {
		opts.color = col
	}
}

func WithPosition(x, y float64) ImageOption {
	return func(opts *ImageOptions) {
		opts.position = geom.Vec2{X: x, Y: y}
	}
}

func WithScale(x, y float64) ImageOption {
	return func(opts *ImageOptions) {
		opts.scale = geom.Vec2{X: x, Y: y}
	}
}

func WithVisible(visible bool) ImageOption {
	return func(opts *ImageOptions) {
		opts.visible = visible
	}
}
