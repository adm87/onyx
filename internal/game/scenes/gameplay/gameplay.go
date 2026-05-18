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

func NewScene(onyx engine.Game) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		SceneID: SceneID,
		OnEnter: func(_ engine.Scene) error {
			return enterScene(onyx.Logger())
		},
		OnUpdate: func(_ engine.Scene, _ float64) (engine.SceneExitCode, error) {
			return updateInput()
		},
		OnDraw: func(scene engine.Scene, screen *ebiten.Image) error {
			if err := scene.Render(screen); err != nil {
				return err
			}

			safeArea := onyx.Screen().SafeArea()
			return renderScene(screen, safeArea)
		},
	}
}

func enterScene(logger engine.Logger) error {
	logger.Info("Entering Gameplay Scene")
	return nil
}

func updateInput() (engine.SceneExitCode, error) {
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
