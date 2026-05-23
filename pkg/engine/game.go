package engine

import (
	"context"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game interface {
	Start() error
	WithContext(ctx context.Context) Game

	Assets() Assets
	Camera() Camera
	Logger() Logger
	Renderer() Renderer
	Scenes() Scenes
	Screen() Screen
	Time() Time
}

type game struct {
	ctx context.Context

	assets   *assets
	camera   *camera
	logger   *logger
	renderer *renderer
	scenes   *scenes
	screen   *screen
	time     *time
}

func setupWindow(title string, width, height int) {
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(width, height)
}

func NewGame(opts ...Option) Game {
	cfg := applyOptions(opts...)

	setupWindow(cfg.Title, cfg.Width, cfg.Height)

	logger := newLogger(os.Stdout)
	assets := newAssets(
		logger,
	)
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
	time := newTime(
		cfg.FPS,
	)
	renderer := newRenderer(
		logger,
	)
	camera := newCamera(
		scenes.world,
		screen,
	)

	return &game{
		ctx:      context.Background(),
		assets:   assets,
		camera:   camera,
		logger:   logger,
		renderer: renderer,
		screen:   screen,
		scenes:   scenes,
		time:     time,
	}
}

func (g *game) Assets() Assets {
	return g.assets
}

func (g *game) Camera() Camera {
	return g.camera
}

func (g *game) Logger() Logger {
	return g.logger
}

func (g *game) Renderer() Renderer {
	return g.renderer
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

		if err := g.renderer.render(g.scenes.world, g.screen.buffer, g.camera.view()); err != nil {
			g.logger.Error("Failed to render scene: %v", err)
			return
		}

		screen.DrawImage(g.screen.buffer, g.screen.options)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screen.Layout(outsideWidth, outsideHeight)
}
