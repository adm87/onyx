package engine

import (
	"os"

	"github.com/adm87/onyx/pkg/indexing/spatialhash"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Game interface {
	Assets() Assets
	Logger() Logger
	Scenes() Scenes
	Screen() Screen
	Time() Time

	Start() error
}

type game struct {
	assets Assets
	input  Input
	logger Logger
	scenes Scenes
	screen Screen
	time   Time
}

func NewGame(width, height int) Game {
	logger := NewLogger(os.Stdout)
	return &game{
		logger: logger,
		assets: NewAssets(logger),
		input:  NewInput(logger),
		scenes: NewScenes(
			spatialhash.New[donburi.Entity](),
			logger,
		),
		screen: NewScreen(width, height, logger),
		time:   NewTime(),
	}
}

func (g *game) Update() error {
	g.Time().Tick()
	return g.scenes.Update(
		g.time.DeltaTime(),
		g.time.FixedDeltaTime(),
		g.time.FixedSteps(),
	)
}

func (g *game) Draw(screen *ebiten.Image) {
	buffer := g.screen.Buffer()
	buffer.Clear()

	g.scenes.Draw(buffer)
	screen.DrawImage(buffer, g.screen.Options())
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screen.Layout(outsideWidth, outsideHeight)
}

func (g *game) Assets() Assets {
	return g.assets
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

func (g *game) Start() error {
	return ebiten.RunGame(g)
}
