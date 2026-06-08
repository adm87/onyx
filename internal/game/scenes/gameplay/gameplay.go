package gameplay

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

const (
	PlayerIdleAnim = "Idle"
	PlayerRunAnim  = "Run"
)

func New() engine.SceneState {

	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			return engine.SceneExitNone, nil
		},
		OnFixedUpdate: func(ecs donburi.World, dt float64) error {
			return nil
		},
		OnLateUpdate: func(ecs donburi.World, dt float64) error {
			return nil
		},
		OnRender: func(ecs donburi.World, img *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error {
			return nil
		},
	}
}
