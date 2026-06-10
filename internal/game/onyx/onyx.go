package onyx

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
)

const (
	SplashScreenSceneID engine.SceneID = "splash_screen"
	GameplaySceneID     engine.SceneID = "gameplay"
)

const (
	SplashScreenCompleteExitCode engine.SceneExitCode = iota + 1
)

type Onyx struct {
	game engine.Game

	images *images.ImageModule
	tiled  *tiled.TiledModule
}

func NewGame(
	game engine.Game,
	images *images.ImageModule,
	tiled *tiled.TiledModule) *Onyx {

	o := &Onyx{
		game:   game,
		images: images,
		tiled:  tiled,
	}

	o.AddScenes()

	return o
}

func (o *Onyx) Start() error {
	return o.game.Start()
}

func (o *Onyx) AddScenes() {
	s := o.game.Scenes()

	s.AddScene(
		SplashScreenSceneID,
		o.SplashScreenScene(),
		engine.SceneTransitions{
			SplashScreenCompleteExitCode: GameplaySceneID,
		},
	)

	s.AddScene(
		GameplaySceneID,
		o.GameplayScene(),
		engine.SceneTransitions{},
	)
}
