package gameplay

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SceneID engine.SceneID = "gameplay"
)

func New(screen engine.Screen, logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		SceneID: SceneID,
		OnEnter: func(scene engine.Scene) error {
			return enterScene(scene, logger)
		},
		OnUpdate: updateInput,
		OnDraw: func(scene engine.Scene, img *ebiten.Image) error {
			if err := scene.Render(img); err != nil {
				return err
			}
			return renderScene(img, screen.SafeArea())
		},
	}
}

func enterScene(scene engine.Scene, logger engine.Logger) error {
	logger.Info("Entering Gameplay Scene")
	_ = scene
	return nil
}

func updateInput(scene engine.Scene, deltaTime float64) (engine.SceneExitCode, error) {
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
