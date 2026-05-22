package splashscreen

import (
	"context"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
)

const CompleteExitCode engine.SceneExitCode = iota

func New(time engine.Time, logger engine.Logger) engine.SceneState {
	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			logger.Info("Entering Splash Screen Scene")
			return nil
		},
		OnExit: func(ctx context.Context, world donburi.World) error {
			logger.Info("Exiting Splash Screen Scene")
			return nil
		},
	}
}
