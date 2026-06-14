package images

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Options struct {
	Handle uint64
	Frame  int
	Anchor geom.Vec2
	Filter ebiten.Filter
	Color  color.RGBA
}

type Option func(*Options)

type ImageModel struct {
	Handle uint64
	Frame  int
	Anchor geom.Vec2
	Filter ebiten.Filter
	Color  color.RGBA
}

var Image = donburi.NewComponentType[ImageModel]()

func defaultImageOptions() *Options {
	return &Options{
		Handle: 0,
		Frame:  0,
		Anchor: geom.Vec2{X: 0, Y: 0},
		Filter: ebiten.FilterNearest,
		Color:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}
}

func WithHandle(handle uint64) Option {
	return func(opts *Options) {
		opts.Handle = handle
	}
}

func WithFrame(frame int) Option {
	return func(opts *Options) {
		opts.Frame = frame
	}
}

func WithAnchor(x, y float64) Option {
	return func(opts *Options) {
		opts.Anchor = geom.Vec2{X: x, Y: y}
	}
}

func WithFilter(filter ebiten.Filter) Option {
	return func(opts *Options) {
		opts.Filter = filter
	}
}

func WithColor(color color.RGBA) Option {
	return func(opts *Options) {
		opts.Color = color
	}
}

func GetImage(entry *donburi.Entry) *ImageModel {
	if !entry.HasComponent(Image) {
		return nil
	}
	return Image.Get(entry)
}

func SetImage(entry *donburi.Entry, options ...Option) {
	opts := defaultImageOptions()
	for _, option := range options {
		option(opts)
	}
	donburi.Add(entry, Image, &ImageModel{
		Handle: opts.Handle,
		Frame:  opts.Frame,
		Anchor: opts.Anchor,
		Filter: opts.Filter,
		Color:  opts.Color,
	})
}

func GetFrame(entry *donburi.Entry) int {
	if !entry.HasComponent(Image) {
		return 0
	}
	return Image.Get(entry).Frame
}

func SetFrame(entry *donburi.Entry, frame int) {
	if !entry.HasComponent(Image) {
		return
	}
	img := Image.Get(entry)
	img.Frame = frame
}

func GetHandle(entry *donburi.Entry) uint64 {
	if !entry.HasComponent(Image) {
		return 0
	}
	return Image.Get(entry).Handle
}

func SetHandle(entry *donburi.Entry, handle uint64) {
	if !entry.HasComponent(Image) {
		return
	}
	img := Image.Get(entry)
	img.Handle = handle
}

func GetAnchor(entry *donburi.Entry) geom.Vec2 {
	if !entry.HasComponent(Image) {
		return defaultImageOptions().Anchor
	}
	return Image.Get(entry).Anchor
}

func SetAnchor(entry *donburi.Entry, x, y float64) {
	if !entry.HasComponent(Image) {
		return
	}
	img := Image.Get(entry)
	img.Anchor = geom.Vec2{X: x, Y: y}
}

func GetFilter(entry *donburi.Entry) ebiten.Filter {
	if !entry.HasComponent(Image) {
		return defaultImageOptions().Filter
	}
	return Image.Get(entry).Filter
}

func SetFilter(entry *donburi.Entry, filter ebiten.Filter) {
	if !entry.HasComponent(Image) {
		return
	}
	img := Image.Get(entry)
	img.Filter = filter
}

func GetColor(entry *donburi.Entry) color.RGBA {
	if !entry.HasComponent(Image) {
		return defaultImageOptions().Color
	}
	return Image.Get(entry).Color
}

func SetColor(entry *donburi.Entry, color color.RGBA) {
	if !entry.HasComponent(Image) {
		return
	}
	img := Image.Get(entry)
	img.Color = color
}

func GetAlpha(entry *donburi.Entry) uint8 {
	if !entry.HasComponent(Image) {
		return defaultImageOptions().Color.A
	}
	return Image.Get(entry).Color.A
}

func SetAlpha(entry *donburi.Entry, alpha uint8) {
	if !entry.HasComponent(Image) {
		return
	}
	img := Image.Get(entry)
	img.Color.A = alpha
}
