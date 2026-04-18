package app

import (
	"image/color"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	name       string
	version    string
	fullscreen bool

	logOutput io.Writer
	logLevel  LogLevel

	screenWidth      int
	screenHeight     int
	screenClearColor color.RGBA
	screenFilter     ebiten.Filter
	screenResizeMode ScreenResizeMode
}

type Option func(*Config)

func WithName(name string) Option {
	return func(cfg *Config) {
		cfg.name = name
	}
}

func WithVersion(version string) Option {
	return func(cfg *Config) {
		cfg.version = version
	}
}

func WithScreenSize(width, height int) Option {
	return func(cfg *Config) {
		cfg.screenWidth = width
		cfg.screenHeight = height
	}
}

func WithFullscreen(fullscreen bool) Option {
	return func(cfg *Config) {
		cfg.fullscreen = fullscreen
	}
}

func WithLogOutput(output io.Writer) Option {
	return func(cfg *Config) {
		cfg.logOutput = output
	}
}

func WithLogLevel(level LogLevel) Option {
	return func(cfg *Config) {
		cfg.logLevel = level
	}
}

func WithScreenClearColor(c color.RGBA) Option {
	return func(cfg *Config) {
		cfg.screenClearColor = c
	}
}

func WithScreenFilter(filter ebiten.Filter) Option {
	return func(cfg *Config) {
		cfg.screenFilter = filter
	}
}

func WithScreenResizeMode(mode ScreenResizeMode) Option {
	return func(cfg *Config) {
		cfg.screenResizeMode = mode
	}
}

func NewConfig(opts ...Option) *Config {
	cfg := &Config{
		name:             "Untitled",
		version:          "",
		screenWidth:      800,
		screenHeight:     600,
		fullscreen:       false,
		logOutput:        io.Discard,
		logLevel:         LogLevelInfo,
		screenClearColor: color.RGBA{0, 0, 0, 255},
		screenFilter:     ebiten.FilterNearest,
		screenResizeMode: ScreenResizeByWidth,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
