package engine

import (
	"context"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game interface {
	Start() error
	WithContext(ctx context.Context) Game

	Logger() Logger
	Scenes() Scenes
	Screen() Screen
	Time() Time
}

type game struct {
	ctx context.Context

	logger *logger
	scenes *scenes
	screen *screen
	time   *time
}

func setupWindow(title string, width, height int) {
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(width, height)
}

func NewGame(opts ...Option) Game {
	cfg := applyOptions(opts...)

	setupWindow(cfg.Title, cfg.Width, cfg.Height)

	logger := newLogger(os.Stdout)
	screen := newScreen(
		cfg.Width,
		cfg.Height,
		cfg.ScreenScale,
		cfg.Filter,
		cfg.BackgroundColor,
		logger,
	)
	scenes := newScenes(
		cfg.InitialScene,
		logger,
	)
	time := newTime(cfg.FPS)

	return &game{
		ctx:    context.Background(),
		logger: logger,
		screen: screen,
		scenes: scenes,
		time:   time,
	}
}

func (g *game) Logger() Logger {
	return g.logger
}

func (g *game) Scenes() Scenes {
	return g.scenes
}

func (g *game) Screen() Screen {
	return g.screen
}

func (g *game) Time() Time {
	return g.time
}

func (g *game) WithContext(ctx context.Context) Game {
	if g.ctx == nil {
		return g
	}
	g.ctx = ctx
	return g
}

func (g *game) Start() error {
	return ebiten.RunGame(g)
}

func (g *game) Update() error {
	select {
	case <-g.ctx.Done():
		return g.ctx.Err()
	default:
		g.time.tick()
		return g.scenes.update(g.ctx)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	select {
	case <-g.ctx.Done():
		return
	default:
		g.screen.buffer.Fill(g.screen.backgroundColor)

		screen.DrawImage(g.screen.buffer, g.screen.options)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screen.Layout(outsideWidth, outsideHeight)
}
