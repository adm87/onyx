package onyx

import (
	"github.com/adm87/onyx/internal/game/scenes/gameplay"
	"github.com/adm87/onyx/internal/game/scenes/splashscreen"
	"github.com/adm87/onyx/pkg/engine"
)

const (
	SplashScreenSceneID engine.SceneID = "splash_screen"
	GameplaySceneID     engine.SceneID = "gameplay"
)

type Onyx struct {
	game engine.Game
}

func NewGame(game engine.Game) *Onyx {
	s := game.Scenes()

	s.AddScene(
		SplashScreenSceneID,
		func() engine.Scene {
			return splashscreen.NewScene(game)
		},
		engine.SceneTransitions{
			splashscreen.SplashScreenCompleteExitCode: GameplaySceneID,
		},
	)
	s.AddScene(
		GameplaySceneID,
		func() engine.Scene {
			return gameplay.NewScene(game)
		},
		engine.SceneTransitions{},
	)

	return &Onyx{
		game: game,
	}
}

func (o *Onyx) Start() error {
	return o.game.Start()
}
