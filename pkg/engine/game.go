package engine

import (
	"context"
	"os"

	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Game interface {
	Start() error
	WithContext(ctx context.Context) Game

	ECS() donburi.World
	Assets() Assets
	Camera() Camera
	Logger() Logger
	Renderer() Renderer
	Scenes() Scenes
	Screen() Screen
	Time() Time
	World() World
}

type game struct {
	ctx context.Context
	ecs donburi.World

	assets   *assets
	camera   *camera
	logger   *logger
	renderer *renderer
	scenes   *scenes
	screen   *screen
	time     *time
	world    *world

	renderables []*donburi.Entry
}

func setupWindow(title string, width, height int) {
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(width, height)
}

func NewGame(opts ...Option) Game {
	cfg := applyOptions(opts...)

	setupWindow(cfg.Title, cfg.Width, cfg.Height)

	ecs := donburi.NewWorld()

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

	renderer := newRenderer(
		logger,
	)

	world := newWorld()

	scenes := newScenes(
		cfg.InitialScene,
		renderer,
		logger,
	)

	time := newTime(
		cfg.FPS,
	)

	camera := newCamera(
		ecs,
		screen,
		transform.NewTransform(ecs).Entity(),
	)

	return &game{
		ctx:         context.Background(),
		ecs:         ecs,
		assets:      assets,
		camera:      camera,
		logger:      logger,
		renderer:    renderer,
		screen:      screen,
		scenes:      scenes,
		time:        time,
		world:       world,
		renderables: make([]*donburi.Entry, 0, 100),
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

func (g *game) World() World {
	return g.world
}

func (g *game) ECS() donburi.World {
	return g.ecs
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
		return g.scenes.update(
			g.ecs,
			g.time.fixedSteps,
			g.time.deltaTime.Seconds(),
			g.time.fixedDeltaTime.Seconds(),
		)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	select {
	case <-g.ctx.Done():
		return
	default:
		g.screen.buffer.Fill(g.screen.backgroundColor)

		viewMatrix := g.camera.View()
		viewport := g.camera.Viewport()

		g.renderables = g.world.QueryInto(g.ecs, viewport, g.renderables[:0])
		g.renderer.render(g.renderables, g.screen.buffer, viewport, viewMatrix)
		g.scenes.render(g.renderables, g.screen.buffer, viewport, viewMatrix)

		screen.DrawImage(g.screen.buffer, g.screen.options)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screen.Layout(outsideWidth, outsideHeight)
}
