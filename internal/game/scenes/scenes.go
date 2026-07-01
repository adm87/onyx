package scenes

import (
	"github.com/adm87/onyx/internal/game/scenes/gameplay"
	"github.com/adm87/onyx/internal/game/scenes/splashscreen"
	"github.com/adm87/onyx/pkg/engine"
)

const (
	SplashScreenSceneID engine.SceneID = "splash_screen"
	GameplaySceneID     engine.SceneID = "gameplay"
)

func RegisterScenes(game engine.Game) {
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
}
