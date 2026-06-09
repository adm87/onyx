package splashscreen

import (
	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
)

const CompleteExitCode engine.SceneExitCode = iota + 1

func New(assets engine.Assets) engine.SceneState {
	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			if err := assets.Load(content.EmbeddedFS(), content.EmbeddedSplash1920x1080Black); err != nil {
				return err
			}
			return nil
		},
		OnExit: func(ecs donburi.World) error {
			assets.Unload(content.EmbeddedSplash1920x1080Black)
			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			return CompleteExitCode, nil
		},
	}
}
