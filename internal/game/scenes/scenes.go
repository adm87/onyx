package scenes

import (
	"github.com/adm87/onyx-game/internal/game/scenes/gameplay"
	"github.com/adm87/onyx-game/internal/game/scenes/splashscreen"
	"github.com/adm87/onyx-game/pkg/engine"
)

const (
	SplashScreenSceneID engine.SceneID = "splash_screen"
	GameplaySceneID     engine.SceneID = "gameplay"
)

func AddScenes(onyx engine.Game) {
	s := onyx.Scenes()

	assets := onyx.Assets()
	time := onyx.Time()
	screen := onyx.Screen()
	camera := onyx.Camera()

	s.AddScene(
		SplashScreenSceneID,
		splashscreen.New(assets, time, screen),
		engine.SceneTransitions{
			splashscreen.CompleteExitCode: GameplaySceneID,
		},
	)

	s.AddScene(
		GameplaySceneID,
		gameplay.New(assets, camera, screen, time),
		engine.SceneTransitions{},
	)
}
