package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	Title string

	Width  int
	Height int

	ScreenScale ScreenScaleMode
	Filter      ebiten.Filter

	BackgroundColor color.RGBA
}

type Option func(*Config)

func WithFilter(filter ebiten.Filter) Option {
	return func(c *Config) {
		c.Filter = filter
	}
}

func WithScreenScale(mode ScreenScaleMode) Option {
	return func(c *Config) {
		c.ScreenScale = mode
	}
}

func WithTitle(title string) Option {
	return func(c *Config) {
		c.Title = title
	}
}

func WithScreenSize(width, height int) Option {
	return func(c *Config) {
		c.Width = width
		c.Height = height
	}
}

func WithBackgroundColor(color color.RGBA) Option {
	return func(c *Config) {
		c.BackgroundColor = color
	}
}

func defaultConfig() *Config {
	return &Config{
		Title:       "Untitled",
		Width:       800,
		Height:      600,
		ScreenScale: ScreenScaleNone,
		Filter:      ebiten.FilterLinear,
		BackgroundColor: color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		},
	}
}

func applyOptions(opts ...Option) *Config {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
