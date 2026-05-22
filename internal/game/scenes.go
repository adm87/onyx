package game

import (
	"github.com/adm87/onyx/internal/game/scenes/splashscreen"
	"github.com/adm87/onyx/pkg/engine"
)

const (
	GameSceneIDSplashScreen engine.SceneID = "splash_screen"
	GameSceneIDGameplay     engine.SceneID = "gameplay"
)

var (
	splashScreenTransitions = engine.SceneTransitions{
		splashscreen.CompleteExitCode: GameSceneIDGameplay,
	}
)

func addScenes(onyx engine.Game) {
	scenes := onyx.Scenes()

	scenes.AddScene(
		GameSceneIDSplashScreen,
		splashscreen.New(
			onyx.Camera(),
			onyx.Time(),
			onyx.Logger(),
		),
		splashScreenTransitions,
	)
}
