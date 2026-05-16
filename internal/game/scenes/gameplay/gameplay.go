package gameplay

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
)

func New(logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		OnEnter: func(w donburi.World) error {
			return enterScene(w, logger)
		},
	}
}

func enterScene(world donburi.World, logger engine.Logger) error {
	logger.Info("Entering Gameplay Scene")
	return nil
}
