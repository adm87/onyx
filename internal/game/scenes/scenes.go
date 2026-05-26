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

func AddScenes(onyx engine.Game) {
	s := onyx.Scenes()

	s.AddScene(
		SplashScreenSceneID,
		splashscreen.New(
			onyx.Assets(),
			onyx.Time(),
			onyx.Screen(),
		),
		engine.SceneTransitions{
			splashscreen.CompleteExitCode: GameplaySceneID,
		},
	)

	s.AddScene(
		GameplaySceneID,
		gameplay.New(
			onyx.Assets(),
			onyx.Camera(),
			onyx.Screen(),
			onyx.Time(),
		),
		engine.SceneTransitions{},
	)
}
