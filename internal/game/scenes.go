package game

import (
	"github.com/adm87/onyx/internal/game/scenes/gameplay"
	"github.com/adm87/onyx/internal/game/scenes/splashscreen"
	"github.com/adm87/onyx/pkg/engine"
)

const (
	SplashScreenSceneID engine.SceneID = "splash_screen"
	GameplaySceneID     engine.SceneID = "gameplay"
)

var (
	splashScreenTransitions = engine.SceneTransitions{
		splashscreen.CompleteExitCode: GameplaySceneID,
	}
)

func addScenes(onyx engine.Game) {
	scenes := onyx.Scenes()

	scenes.AddScene(
		SplashScreenSceneID,
		splashscreen.New(
			onyx.Assets(),
			onyx.Time(),
			onyx.Screen(),
		),
		splashScreenTransitions,
	)

	scenes.AddScene(
		GameplaySceneID,
		gameplay.New(
			onyx.Assets(),
			onyx.Camera(),
			onyx.Time(),
		),
		engine.SceneTransitions{},
	)
}
