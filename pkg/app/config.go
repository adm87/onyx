package app

import (
	"image/color"
	"io"
)

type Config struct {
	name       string
	version    string
	width      int
	height     int
	fullscreen bool

	logOutput io.Writer
	logLevel  LogLevel

	screenClearColor color.RGBA
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
		cfg.width = width
		cfg.height = height
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

func NewConfig(opts ...Option) *Config {
	cfg := &Config{
		name:             "Untitled",
		version:          "",
		width:            800,
		height:           600,
		fullscreen:       false,
		logOutput:        io.Discard,
		logLevel:         LogLevelInfo,
		screenClearColor: color.RGBA{0, 0, 0, 255},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
