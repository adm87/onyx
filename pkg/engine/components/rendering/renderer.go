package rendering

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type rendererInfo struct {
	Visible bool
	Layer   int
	ZIndex  int
}

type (
	Options struct {
		Visible bool
		Layer   int
		ZIndex  int
		Anchor  geom.Vec2
		Color   color.RGBA
		Filter  ebiten.Filter
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
)

var (
	defaultFilter   = ebiten.FilterNearest
	defaultColor    = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	defaultAnchor   = geom.Vec2{X: 0, Y: 0}
	defaultRenderer = rendererInfo{Visible: true}
)

var (
	Filter   = donburi.NewComponentType[ebiten.Filter](defaultFilter)
	Anchor   = donburi.NewComponentType[geom.Vec2](defaultAnchor)
	Color    = donburi.NewComponentType[color.RGBA](defaultColor)
	renderer = donburi.NewComponentType[rendererInfo](defaultRenderer)
)

var query = donburi.NewQuery(
	filter.Contains(
		Filter,
		Anchor,
		Color,
		renderer,
	),
)

func defaultOptions() Options {
	return Options{
		Visible: defaultRenderer.Visible,
		Layer:   defaultRenderer.Layer,
		ZIndex:  defaultRenderer.ZIndex,
		Color:   defaultColor,
		Filter:  defaultFilter,
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

// Query iterates over all entities with the necessary components for rendering and applies the provided function to each entry.
func Query(ecs donburi.World, fn QueryCallback) {
	query.Each(ecs, func(entry *donburi.Entry) {
		anchor := GetAnchor(entry)
		color := GetColor(entry)
		filter := GetFilter(entry)
		info := renderer.Get(entry)
		fn(entry, anchor, color, filter, info.Visible, info.Layer, info.ZIndex)
	})
}

// QueryWith iterates over all entities that match the provided query and have the necessary components for rendering,
// then applies the provided function to each entry.
func QueryWith(ecs donburi.World, q *donburi.Query, fn QueryCallback) {
	q.Each(ecs, func(entry *donburi.Entry) {
		anchor := GetAnchor(entry)
		color := GetColor(entry)
		filter := GetFilter(entry)
		info := renderer.Get(entry)
		fn(entry, anchor, color, filter, info.Visible, info.Layer, info.ZIndex)
	})
}

// QueryVisible iterates over all entities with the necessary components for rendering that are also marked as visible,
func QueryVisible(ecs donburi.World, fn QueryCallback) {
	query.Each(ecs, func(entry *donburi.Entry) {
		info := renderer.Get(entry)
		if !info.Visible {
			return
		}
		anchor := GetAnchor(entry)
		color := GetColor(entry)
		filter := GetFilter(entry)
		fn(entry, anchor, color, filter, info.Visible, info.Layer, info.ZIndex)
	})
}

// QueryVisibleWith iterates over all entities that match the provided query and have the necessary components for rendering,
// then applies the provided function to each entry that is also marked as visible.
func QueryVisibleWith(ecs donburi.World, q *donburi.Query, fn QueryCallback) {
	q.Each(ecs, func(entry *donburi.Entry) {
		info := renderer.Get(entry)
		if !info.Visible {
			return
		}
		anchor := GetAnchor(entry)
		color := GetColor(entry)
		filter := GetFilter(entry)
		fn(entry, anchor, color, filter, info.Visible, info.Layer, info.ZIndex)
	})
}

// NewRenderer creates a new entity with the necessary components for rendering and applies the provided options.
func NewRenderer(ecs donburi.World, options ...Option) *donburi.Entry {
	return AddRenderer(ecs.Entry(
		ecs.Create(
			Filter,
			Anchor,
			Color,
			renderer,
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
	SetAnchor(entry, opts.Anchor)
	SetColor(entry, opts.Color)

	donburi.Add(entry, renderer, &rendererInfo{
		Visible: opts.Visible,
		Layer:   opts.Layer,
		ZIndex:  opts.ZIndex,
	})

	return entry
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
func SetAnchor(entry *donburi.Entry, anchor geom.Vec2) {
	donburi.Add(entry, Anchor, &anchor)
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

// GetLayer retrieves the layer information from an entity, returning a default value if it does not exist.
func GetLayer(entry *donburi.Entry) int {
	if !entry.HasComponent(renderer) {
		return defaultRenderer.Layer
	}
	return renderer.Get(entry).Layer
}

// SetLayer sets the layer information for an entity, adding it if it does not already exist.
func SetLayer(entry *donburi.Entry, layer int) {
	if !entry.HasComponent(renderer) {
		info := defaultRenderer
		info.Layer = layer
		donburi.Add(entry, renderer, &info)
		return
	}
	info := renderer.Get(entry)
	info.Layer = layer
}

// GetZIndex retrieves the z-index information from an entity, returning a default value if it does not exist.
func GetZIndex(entry *donburi.Entry) int {
	if !entry.HasComponent(renderer) {
		return defaultRenderer.ZIndex
	}
	return renderer.Get(entry).ZIndex
}

// SetZIndex sets the z-index information for an entity, adding it if it does not already exist.
func SetZIndex(entry *donburi.Entry, zIndex int) {
	if !entry.HasComponent(renderer) {
		info := defaultRenderer
		info.ZIndex = zIndex
		donburi.Add(entry, renderer, &info)
		return
	}
	info := renderer.Get(entry)
	info.ZIndex = zIndex
}

// IsVisible retrieves the visibility information from an entity, returning a default value if it does not exist.
func IsVisible(entry *donburi.Entry) bool {
	if !entry.HasComponent(renderer) {
		return defaultRenderer.Visible
	}
	return renderer.Get(entry).Visible
}

// SetVisible sets the visibility information for an entity, adding it if it does not already exist.
func SetVisible(entry *donburi.Entry, visible bool) {
	if !entry.HasComponent(renderer) {
		info := defaultRenderer
		info.Visible = visible
		donburi.Add(entry, renderer, &info)
		return
	}
	info := renderer.Get(entry)
	info.Visible = visible
}
