package game

import (
	"context"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Options struct {
	ctx context.Context

	mdl Model
	lgr Logger

	title   string
	version string

	width  int
	height int
	fps    int

	fullscreen bool

	OnStart       GameFunc
	OnUpdate      GameFunc
	OnFixedUpdate GameFunc
	OnLateUpdate  GameFunc
	OnDraw        GameFunc

	clearColor color.RGBA

	screenFilter ebiten.Filter
	runOpts      ebiten.RunGameOptions
}

func defaultOpts() *Options {
	return &Options{
		ctx:           context.Background(),
		mdl:           &modelImpl{title: "Untitled", version: "1.0.0"},
		lgr:           &loggerImpl{out: os.Stdout},
		title:         "Untitled",
		version:       "1.0.0",
		width:         800,
		height:        600,
		fps:           60,
		fullscreen:    false,
		clearColor:    color.RGBA{0, 0, 0, 255},
		OnStart:       gameNoop,
		OnUpdate:      gameNoop,
		OnFixedUpdate: gameNoop,
		OnLateUpdate:  gameNoop,
		OnDraw:        gameNoop,
		screenFilter:  ebiten.FilterNearest,
		runOpts:       ebiten.RunGameOptions{},
	}
}

func applyOpts(base *Options, opts ...Option) {
	for _, opt := range opts {
		opt(base)
	}
}

type Option func(*Options)

func WithContext(ctx context.Context) Option {
	if ctx == nil {
		ctx = context.Background()
	}
	return func(opts *Options) {
		opts.ctx = ctx
	}
}

func WithModel(mdl Model) Option {
	if mdl == nil {
		mdl = &modelImpl{title: "Untitled", version: "1.0.0"}
	}
	return func(opts *Options) {
		opts.mdl = mdl
	}
}

func WithLogger(lgr Logger) Option {
	if lgr == nil {
		lgr = &loggerImpl{out: os.Stdout}
	}
	return func(opts *Options) {
		opts.lgr = lgr
	}
}

func WithTitle(title string) Option {
	if title == "" {
		title = "Untitled"
	}
	return func(opts *Options) {
		opts.title = title
	}
}

func WithVersion(version string) Option {
	if version == "" {
		version = "1.0.0"
	}
	return func(opts *Options) {
		opts.version = version
	}
}

func WithSize(width, height int) Option {
	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 600
	}
	return func(opts *Options) {
		opts.width = width
		opts.height = height
	}
}

func WithFPS(fps int) Option {
	if fps <= 0 {
		fps = 60
	}
	return func(opts *Options) {
		opts.fps = fps
	}
}

func WithFullscreen(fullscreen bool) Option {
	return func(opts *Options) {
		opts.fullscreen = fullscreen
	}
}

// WithOnStart sets the function to be called when the game starts.
// Called once before the game loop begins.
// If nil is provided, it defaults to a no-op function.
func WithOnStart(onStart GameFunc) Option {
	if onStart == nil {
		onStart = gameNoop
	}
	return func(opts *Options) {
		opts.OnStart = onStart
	}
}

// WithOnUpdate sets the function to be called on each update tick.
// Called once per frame before fixed updates and late updates.
// If nil is provided, it defaults to a no-op function.
func WithOnUpdate(onUpdate GameFunc) Option {
	if onUpdate == nil {
		onUpdate = gameNoop
	}
	return func(opts *Options) {
		opts.OnUpdate = onUpdate
	}
}

// WithOnFixedUpdate sets the function to be called on each fixed update tick.
// Called at a fixed interval, independent of frame rate, after the regular update and before late updates.
// If nil is provided, it defaults to a no-op function.
func WithOnFixedUpdate(onFixedUpdate GameFunc) Option {
	if onFixedUpdate == nil {
		onFixedUpdate = gameNoop
	}
	return func(opts *Options) {
		opts.OnFixedUpdate = onFixedUpdate
	}
}

// WithOnLateUpdate sets the function to be called on each late update tick.
// Called once per frame after regular updates and fixed updates, ideal for cleanup or post-processing.
// If nil is provided, it defaults to a no-op function.
func WithOnLateUpdate(onLateUpdate GameFunc) Option {
	if onLateUpdate == nil {
		onLateUpdate = gameNoop
	}
	return func(opts *Options) {
		opts.OnLateUpdate = onLateUpdate
	}
}

// WithOnDraw sets the function to be called on each draw tick.
// Called once per frame after all updates, responsible for rendering the game state to the screen.
// If nil is provided, it defaults to a no-op function.
func WithOnDraw(onDraw GameFunc) Option {
	if onDraw == nil {
		onDraw = gameNoop
	}
	return func(opts *Options) {
		opts.OnDraw = onDraw
	}
}

func WithClearColor(clearColor color.RGBA) Option {
	return func(opts *Options) {
		opts.clearColor = clearColor
	}
}

func WithScreenFilter(filter ebiten.Filter) Option {
	return func(opts *Options) {
		opts.screenFilter = filter
	}
}

func WithGraphicsLibrary(graphicsLibrary ebiten.GraphicsLibrary) Option {
	return func(opts *Options) {
		opts.runOpts.GraphicsLibrary = graphicsLibrary
	}
}

func WithInitUnfocused(initUnfocused bool) Option {
	return func(opts *Options) {
		opts.runOpts.InitUnfocused = initUnfocused
	}
}

func WithScreenTransparent(screenTransparent bool) Option {
	return func(opts *Options) {
		opts.runOpts.ScreenTransparent = screenTransparent
	}
}

func WithSkipTaskbar(skipTaskbar bool) Option {
	return func(opts *Options) {
		opts.runOpts.SkipTaskbar = skipTaskbar
	}
}

func WithSingleThread(singleThread bool) Option {
	return func(opts *Options) {
		opts.runOpts.SingleThread = singleThread
	}
}

func WithDisableHiDPI(disableHiDPI bool) Option {
	return func(opts *Options) {
		opts.runOpts.DisableHiDPI = disableHiDPI
	}
}

func WithColorSpace(colorSpace ebiten.ColorSpace) Option {
	return func(opts *Options) {
		opts.runOpts.ColorSpace = colorSpace
	}
}

func WithApplePressAndHold(applePressAndHold bool) Option {
	return func(opts *Options) {
		opts.runOpts.ApplePressAndHoldEnabled = applePressAndHold
	}
}

func WithX11ClassName(className string) Option {
	return func(opts *Options) {
		opts.runOpts.X11ClassName = className
	}
}

func WithX11InstanceName(instanceName string) Option {
	return func(opts *Options) {
		opts.runOpts.X11InstanceName = instanceName
	}
}
