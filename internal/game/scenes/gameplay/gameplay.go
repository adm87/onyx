package gameplay

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
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
			onyx.Logger().Info("Entering Gameplay Scene")
			return nil
		},
		OnExit: func(_ engine.Scene) error {
			onyx.Logger().Info("Exiting Gameplay Scene")
			return nil
		},
		OnUpdate: func(_ engine.Scene, _ float64) (engine.SceneExitCode, error) {
			if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
				return engine.SceneExitCodeNone, ebiten.Termination
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF) {
				ebiten.SetFullscreen(!ebiten.IsFullscreen())
			}
			return engine.SceneExitCodeNone, nil
		},
		OnDraw: func(scene engine.Scene, screen *ebiten.Image) error {
			if err := scene.Render(screen, geom.AABB{}, ebiten.GeoM{}); err != nil {
				return err
			}

			safeArea := onyx.Screen().SafeArea()

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
		},
	}
}
