package game

import (
	"github.com/adm87/onyx/internal/game/scenes/gameplay"
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
			onyx.Assets(),
			onyx.Time(),
			onyx.Screen(),
		),
		splashScreenTransitions,
	)

	scenes.AddScene(
		GameSceneIDGameplay,
		gameplay.New(
			onyx.Assets(),
			onyx.Camera(),
			onyx.Time(),
		),
		engine.SceneTransitions{},
	)
}
