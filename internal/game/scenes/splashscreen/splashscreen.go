package splashscreen

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
)

const CompleteExitCode engine.SceneExitCode = iota + 1

func New() engine.SceneState {
	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			return nil
		},
		OnExit: func(ecs donburi.World) error {
			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			return engine.SceneExitNone, nil
		},
	}
}
