package splashscreen

import (
	"time"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SceneID          engine.SceneID       = "splashscreen"
	CompleteExitCode engine.SceneExitCode = iota + 1
)

func NewScene(onyx engine.Game) *engine.SceneDefinition {
	ticker := time.NewTicker(time.Second * 2)
	return &engine.SceneDefinition{
		SceneID: SceneID,
		OnEnter: func(_ engine.Scene) error {
			onyx.Logger().Info("Entering Splash Screen Scene")
			return nil
		},
		OnExit: func(_ engine.Scene) error {
			onyx.Logger().Info("Exiting Splash Screen Scene")
			return nil
		},
		OnUpdate: func(_ engine.Scene, _ float64) (engine.SceneExitCode, error) {
			select {
			case <-ticker.C:
				ticker.Stop()
				return CompleteExitCode, nil
			default:
				return engine.SceneExitCodeNone, nil
			}
		},
		OnDraw: func(scene engine.Scene, screen *ebiten.Image) error {
			return scene.Render(screen, geom.AABB{}, ebiten.GeoM{})
		},
	}
}
