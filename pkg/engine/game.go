package engine

import (
	"context"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game interface {
	Start() error
	WithContext(ctx context.Context) Game

	Logger() Logger
	Screen() Screen
}

type game struct {
	ctx context.Context

	logger *logger
	screen *screen
}

func setupWindow(title string, width, height int) {
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizable(true)
}

func NewGame(opts ...Option) Game {
	cfg := applyOptions(opts...)
	logger := newLogger(os.Stdout)

	setupWindow(cfg.Title, cfg.Width, cfg.Height)

	return &game{
		ctx:    context.Background(),
		logger: logger,
		screen: newScreen(
			cfg.Width,
			cfg.Height,
			cfg.ScreenScale,
			cfg.Filter,
			cfg.BackgroundColor,
			logger,
		),
	}
}

func (s *game) Logger() Logger {
	return s.logger
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
		if inpututil.IsKeyJustPressed(ebiten.KeyF) {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}
		return nil
	}
}

func (s *game) Draw(screen *ebiten.Image) {
	select {
	case <-s.ctx.Done():
		return
	default:
		s.screen.buffer.Fill(s.screen.backgroundColor)

		// Top Left
		vector.FillRect(s.screen.buffer,
			float32(s.screen.safeArea.Min.X+10),
			float32(s.screen.safeArea.Min.Y+10),
			100, 100, color.White, false)

		// Top Right
		vector.FillRect(s.screen.buffer,
			float32(s.screen.safeArea.Max.X-110),
			float32(s.screen.safeArea.Min.Y+10),
			100, 100, color.White, false)

		// Bottom Left
		vector.FillRect(s.screen.buffer,
			float32(s.screen.safeArea.Min.X+10),
			float32(s.screen.safeArea.Max.Y-110),
			100, 100, color.White, false)

		// Bottom Right
		vector.FillRect(s.screen.buffer,
			float32(s.screen.safeArea.Max.X-110),
			float32(s.screen.safeArea.Max.Y-110),
			100, 100, color.White, false)

		screen.DrawImage(s.screen.buffer, s.screen.options)
	}
}

func (s *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.screen.Layout(outsideWidth, outsideHeight)
}
