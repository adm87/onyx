package game

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	assets engine.Assets
	logger engine.Logger
	scenes engine.Scenes
	screen engine.Screen
	time   engine.Time
}

func (g *game) Update() error {
	g.time.Tick()
	g.scenes.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.screen.Clear()

	buffer := g.screen.Buffer()

	ebitenutil.DebugPrint(buffer, "Hello, Onyx!")

	screen.DrawImage(buffer, g.screen.DrawOptions())
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screen.HandleLayout(outsideWidth, outsideHeight)
}
