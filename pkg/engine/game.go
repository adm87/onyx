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
}

type game struct {
	ctx context.Context

	logger *logger
	scenes *scenes
	screen *screen
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

	return &game{
		ctx:    context.Background(),
		logger: logger,
		screen: screen,
		scenes: scenes,
	}
}

func (s *game) Logger() Logger {
	return s.logger
}

func (s *game) Scenes() Scenes {
	return s.scenes
}

func (s *game) Screen() Screen {
	return s.screen
}

func (s *game) WithContext(ctx context.Context) Game {
	if s.ctx == nil {
		return s
	}
	s.ctx = ctx
	return s
}

func (s *game) Start() error {
	return ebiten.RunGame(s)
}

func (s *game) Update() error {
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
		return s.scenes.update(s.ctx)
	}
}

func (s *game) Draw(screen *ebiten.Image) {
	select {
	case <-s.ctx.Done():
		return
	default:
		s.screen.buffer.Fill(s.screen.backgroundColor)

		screen.DrawImage(s.screen.buffer, s.screen.options)
	}
}

func (s *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.screen.Layout(outsideWidth, outsideHeight)
}
