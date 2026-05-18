package splashscreen

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SceneID          engine.SceneID       = "splashscreen"
	CompleteExitCode engine.SceneExitCode = iota + 1
)

func New(assets engine.Assets, screen engine.Screen, time engine.Time, logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		SceneID: SceneID,
		OnEnter: func(_ engine.Scene) error {
			return enterScene(logger)
		},
		OnExit: func(_ engine.Scene) error {
			return exitScene(logger)
		},
		OnDraw: func(s engine.Scene, i *ebiten.Image) error {
			return s.Render(i)
		},
	}
}

func enterScene(_ engine.Logger) error {
	return nil
}

func exitScene(_ engine.Logger) error {
	return nil
}
