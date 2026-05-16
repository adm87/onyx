package scenes

import (
	"github.com/adm87/onyx/internal/game/scenes/gameplay"
	"github.com/adm87/onyx/internal/game/scenes/splashscreen"
	"github.com/adm87/onyx/pkg/engine"
)

const (
	SplashScreenSceneID engine.SceneID = iota + 1
	GameplaySceneID
)

func Register(game engine.Game) {
	gameScenes := game.Scenes()

	// Splash Screen
	gameScenes.Register(SplashScreenSceneID,
		splashscreen.New(
			game.Assets(),
			game.Screen(),
			game.Time(),
			game.Logger(),
		),
		engine.SceneTransitions{
			splashscreen.CompleteExitCode: GameplaySceneID,
		},
	)

	// Gameplay
	gameScenes.Register(GameplaySceneID,
		gameplay.New(
			game.Logger(),
		),
		engine.SceneTransitions{},
	)
}
