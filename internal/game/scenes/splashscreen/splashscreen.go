package splashscreen

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
)

const (
	SceneID          engine.SceneID       = "splashscreen"
	CompleteExitCode engine.SceneExitCode = iota + 1
)

func New(assets engine.Assets, screen engine.Screen, time engine.Time, logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		SceneID: SceneID,
		OnEnter: func(w donburi.World) error {
			return enterScene(logger)
		},
		OnExit: func(w donburi.World) error {
			return exitScene(logger)
		},
	}
}

func enterScene(logger engine.Logger) error {
	return nil
}

func exitScene(logger engine.Logger) error {
	return nil
}
