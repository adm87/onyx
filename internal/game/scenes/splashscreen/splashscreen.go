package splashscreen

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
)

const (
	CompleteExitCode engine.SceneExitCode = iota + 1
)

func New(assets engine.Assets, screen engine.Screen, logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		OnEnter: func(w donburi.World) error {
			return enterScene(w, assets, screen, logger)
		},
		OnExit: func(w donburi.World) error {
			return exitScene(w, logger)
		},
	}
}

func enterScene(world donburi.World, assets engine.Assets, screen engine.Screen, logger engine.Logger) error {
	logger.Info("Entering Splash Screen Scene")
	_ = world
	_ = assets
	_ = screen
	return nil
}

func exitScene(world donburi.World, logger engine.Logger) error {
	logger.Info("Exiting Splash Screen Scene")
	_ = world
	return nil
}
