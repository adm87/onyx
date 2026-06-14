package rendering

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type RendererInfo struct {
	Visible bool
	Layer   int
	ZIndex  int
	ID      uint64
}

type (
	Options struct {
		Visible  bool
		Layer    int
		ZIndex   int
		Renderer uint64
		Anchor   geom.Vec2
		Color    color.RGBA
		Filter   ebiten.Filter
	}
	Option func(*Options)
)

type QueryCallback func(
	world *donburi.Entry,
	anchor geom.Vec2,
	color color.RGBA,
	filter ebiten.Filter,
	visible bool,
	layer int,
	zIndex int,
	renderer uint64,
)

var (
	defaultFilter   = ebiten.FilterNearest
	defaultColor    = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	defaultAnchor   = geom.Vec2{X: 0, Y: 0}
	defaultRenderer = RendererInfo{Visible: true}
)

var (
	Anchor   = donburi.NewComponentType[geom.Vec2](defaultAnchor)
	Color    = donburi.NewComponentType[color.RGBA](defaultColor)
	Filter   = donburi.NewComponentType[ebiten.Filter](defaultFilter)
	Renderer = donburi.NewComponentType[RendererInfo](defaultRenderer)
)

var query = donburi.NewQuery(
	filter.Contains(
		Filter,
		Anchor,
		Color,
		Renderer,
	),
)

func defaultOptions() Options {
	return Options{
		Visible: defaultRenderer.Visible,
		Layer:   defaultRenderer.Layer,
		ZIndex:  defaultRenderer.ZIndex,
		Anchor:  defaultAnchor,
		Color:   defaultColor,
		Filter:  defaultFilter,
	}
}

func WithRendererID(renderer uint64) Option {
	return func(opts *Options) {
		opts.Renderer = renderer
	}
}

func WithVisibility(visible bool) Option {
	return func(opts *Options) {
		opts.Visible = visible
	}
}

func WithLayer(layer int) Option {
	return func(opts *Options) {
		opts.Layer = layer
	}
}

func WithZIndex(zIndex int) Option {
	return func(opts *Options) {
		opts.ZIndex = zIndex
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

func WithFilter(filter ebiten.Filter) Option {
	return func(opts *Options) {
		opts.Filter = filter
	}
}

// NewRenderer creates a new entity with the necessary components for rendering and applies the provided options.
func NewRenderer(ecs donburi.World, options ...Option) *donburi.Entry {
	return AddRenderer(ecs.Entry(
		ecs.Create(
			Filter,
			Anchor,
			Color,
			Renderer,
		),
	), options...)
}

// AddRenderer adds the necessary components for rendering to an existing entity and applies the provided options.
func AddRenderer(entry *donburi.Entry, options ...Option) *donburi.Entry {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}

	SetFilter(entry, opts.Filter)
	SetAnchor(entry, opts.Anchor.X, opts.Anchor.Y)
	SetColor(entry, opts.Color)

	donburi.Add(entry, Renderer, &RendererInfo{
		Visible: opts.Visible,
		Layer:   opts.Layer,
		ZIndex:  opts.ZIndex,
		ID:      opts.Renderer,
	})

	return entry
}

func HasRenderer(entry *donburi.Entry) bool {
	return entry.HasComponent(Renderer)
}

func GetRenderer(entry *donburi.Entry) uint64 {
	if !entry.HasComponent(Renderer) {
		return defaultRenderer.ID
	}
	return Renderer.Get(entry).ID
}

// GetFilter retrieves the filter component from an entity, returning a default value if it does not exist.
func GetFilter(entry *donburi.Entry) ebiten.Filter {
	if !entry.HasComponent(Filter) {
		return defaultFilter
	}
	return *Filter.Get(entry)
}

// SetFilter sets the filter component for an entity, adding it if it does not already exist.
func SetFilter(entry *donburi.Entry, filter ebiten.Filter) {
	donburi.Add(entry, Filter, &filter)
}

// GetAnchor retrieves the anchor component from an entity, returning a default value if it does not exist.
func GetAnchor(entry *donburi.Entry) geom.Vec2 {
	if !entry.HasComponent(Anchor) {
		return defaultAnchor
	}
	return *Anchor.Get(entry)
}

// SetAnchor sets the anchor component for an entity, adding it if it does not already exist.
func SetAnchor(entry *donburi.Entry, x, y float64) {
	donburi.Add(entry, Anchor, &geom.Vec2{X: x, Y: y})
}

// GetColor retrieves the color component from an entity, returning a default value if it does not exist.
func GetColor(entry *donburi.Entry) color.RGBA {
	if !entry.HasComponent(Color) {
		return defaultColor
	}
	return *Color.Get(entry)
}

// SetColor sets the color component for an entity, adding it if it does not already exist.
func SetColor(entry *donburi.Entry, color color.RGBA) {
	donburi.Add(entry, Color, &color)
}

func SetAlpha(entry *donburi.Entry, alpha uint8) {
	color := GetColor(entry)
	color.A = alpha
	SetColor(entry, color)
}

// GetLayer retrieves the layer information from an entity, returning a default value if it does not exist.
func GetLayer(entry *donburi.Entry) int {
	if !entry.HasComponent(Renderer) {
		return defaultRenderer.Layer
	}
	return Renderer.Get(entry).Layer
}

// SetLayer sets the layer information for an entity, adding it if it does not already exist.
func SetLayer(entry *donburi.Entry, layer int) {
	if !entry.HasComponent(Renderer) {
		info := defaultRenderer
		info.Layer = layer
		donburi.Add(entry, Renderer, &info)
		return
	}
	info := Renderer.Get(entry)
	info.Layer = layer
}

// GetZIndex retrieves the z-index information from an entity, returning a default value if it does not exist.
func GetZIndex(entry *donburi.Entry) int {
	if !entry.HasComponent(Renderer) {
		return defaultRenderer.ZIndex
	}
	return Renderer.Get(entry).ZIndex
}

// SetZIndex sets the z-index information for an entity, adding it if it does not already exist.
func SetZIndex(entry *donburi.Entry, zIndex int) {
	if !entry.HasComponent(Renderer) {
		info := defaultRenderer
		info.ZIndex = zIndex
		donburi.Add(entry, Renderer, &info)
		return
	}
	info := Renderer.Get(entry)
	info.ZIndex = zIndex
}

// IsVisible retrieves the visibility information from an entity, returning a default value if it does not exist.
func IsVisible(entry *donburi.Entry) bool {
	if !entry.HasComponent(Renderer) {
		return defaultRenderer.Visible
	}
	return Renderer.Get(entry).Visible
}

// SetVisible sets the visibility information for an entity, adding it if it does not already exist.
func SetVisible(entry *donburi.Entry, visible bool) {
	if !entry.HasComponent(Renderer) {
		info := defaultRenderer
		info.Visible = visible
		donburi.Add(entry, Renderer, &info)
		return
	}
	info := Renderer.Get(entry)
	info.Visible = visible
}
