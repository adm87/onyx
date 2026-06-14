package rendering

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RendererModel struct {
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

var (
	defaultFilter   = ebiten.FilterNearest
	defaultColor    = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	defaultAnchor   = geom.Vec2{X: 0, Y: 0}
	defaultRenderer = RendererModel{Visible: true}
)

var (
	RendererID = donburi.NewComponentType[RendererModel](defaultRenderer)
	Renderer   = donburi.NewComponentType[RendererModel]()
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
			RendererID,
		),
	), options...)
}

// AddRenderer adds the necessary components for rendering to an existing entity and applies the provided options.
func AddRenderer(entry *donburi.Entry, options ...Option) *donburi.Entry {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}

	donburi.Add(entry, RendererID, &RendererModel{
		Visible: opts.Visible,
		Layer:   opts.Layer,
		ZIndex:  opts.ZIndex,
		ID:      opts.Renderer,
	})

	return entry
}

func HasRenderer(entry *donburi.Entry) bool {
	return entry.HasComponent(RendererID)
}

func GetRendererID(entry *donburi.Entry) uint64 {
	if !entry.HasComponent(RendererID) {
		return defaultRenderer.ID
	}
	return RendererID.Get(entry).ID
}

// GetLayer retrieves the layer information from an entity, returning a default value if it does not exist.
func GetLayer(entry *donburi.Entry) int {
	if !entry.HasComponent(RendererID) {
		return defaultRenderer.Layer
	}
	return RendererID.Get(entry).Layer
}

// SetLayer sets the layer information for an entity, adding it if it does not already exist.
func SetLayer(entry *donburi.Entry, layer int) {
	if !entry.HasComponent(RendererID) {
		info := defaultRenderer
		info.Layer = layer
		donburi.Add(entry, RendererID, &info)
		return
	}
	info := RendererID.Get(entry)
	info.Layer = layer
}

// GetZIndex retrieves the z-index information from an entity, returning a default value if it does not exist.
func GetZIndex(entry *donburi.Entry) int {
	if !entry.HasComponent(RendererID) {
		return defaultRenderer.ZIndex
	}
	return RendererID.Get(entry).ZIndex
}

// SetZIndex sets the z-index information for an entity, adding it if it does not already exist.
func SetZIndex(entry *donburi.Entry, zIndex int) {
	if !entry.HasComponent(RendererID) {
		info := defaultRenderer
		info.ZIndex = zIndex
		donburi.Add(entry, RendererID, &info)
		return
	}
	info := RendererID.Get(entry)
	info.ZIndex = zIndex
}

// IsVisible retrieves the visibility information from an entity, returning a default value if it does not exist.
func IsVisible(entry *donburi.Entry) bool {
	if !entry.HasComponent(RendererID) {
		return defaultRenderer.Visible
	}
	return RendererID.Get(entry).Visible
}

// SetVisible sets the visibility information for an entity, adding it if it does not already exist.
func SetVisible(entry *donburi.Entry, visible bool) {
	if !entry.HasComponent(RendererID) {
		info := defaultRenderer
		info.Visible = visible
		donburi.Add(entry, RendererID, &info)
		return
	}
	info := RendererID.Get(entry)
	info.Visible = visible
}

func GetRenderer(entry *donburi.Entry) *RendererModel {
	if !entry.HasComponent(RendererID) {
		return &defaultRenderer
	}
	return RendererID.Get(entry)
}

func SetRenderer(entry *donburi.Entry, options ...Option) {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	donburi.Add(entry, RendererID, &RendererModel{
		Visible: opts.Visible,
		Layer:   opts.Layer,
		ZIndex:  opts.ZIndex,
		ID:      opts.Renderer,
	})
}
