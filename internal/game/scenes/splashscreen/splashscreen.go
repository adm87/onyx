package splashscreen

import (
	"github.com/adm87/onyx/pkg/engine"
)

const (
	CompleteExitCode engine.SceneExitCode = iota + 1
)

func New(assets engine.Assets, screen engine.Screen, time engine.Time, logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{}
}
