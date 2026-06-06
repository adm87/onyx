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

	assets := onyx.Assets()
	camera := onyx.Camera()
	collision := onyx.Collision()
	screen := onyx.Screen()
	time := onyx.Time()
	world := onyx.World()

	s.AddScene(
		SplashScreenSceneID,
		splashscreen.New(
			assets,
			time,
			screen,
		),
		engine.SceneTransitions{
			splashscreen.CompleteExitCode: GameplaySceneID,
		},
	)

	s.AddScene(
		GameplaySceneID,
		gameplay.New(
			assets,
			camera,
			collision,
			screen,
			world,
			time,
		),
		engine.SceneTransitions{},
	)
}
