package gameplay

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

func New(screen engine.Screen, logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		OnEnter: func(w donburi.World) error {
			return enterScene(w, logger)
		},
		OnUpdate: updateInput,
		OnDraw: func(_ donburi.World, img *ebiten.Image) error {
			return renderScene(img, screen.SafeArea())
		},
	}
}

func enterScene(_ donburi.World, logger engine.Logger) error {
	logger.Info("Entering Gameplay Scene")
	return nil
}

func updateInput(world donburi.World, deltaTime float64) (engine.SceneExitCode, error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return engine.SceneExitCodeNone, ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return engine.SceneExitCodeNone, nil
}

func renderScene(screen *ebiten.Image, safeArea *engine.SafeArea) error {
	minX, minY := safeArea.Min()
	maxX, maxY := safeArea.Max()

	left := float32(minX + 10)
	top := float32(minY + 10)
	right := float32(maxX - 10)
	bottom := float32(maxY - 10)

	vector.FillRect(screen, left, top, 100, 100, color.White, false)
	vector.FillRect(screen, right-100, top, 100, 100, color.White, false)
	vector.FillRect(screen, left, bottom-100, 100, 100, color.White, false)
	vector.FillRect(screen, right-100, bottom-100, 100, 100, color.White, false)
	return nil
}
