package onyx

import (
	"github.com/adm87/onyx/pkg/ecs"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/aseprite"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/adm87/onyx/pkg/plugins/tiled"
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

	animations *aseprite.AsepritePlugin
	ecs        *ecs.DonburiECSPlugin
	images     *images.ImagePlugin
	tiled      *tiled.TiledPlugin
}

func NewGame(
	game engine.Game,
	animations *aseprite.AsepritePlugin,
	ecs *ecs.DonburiECSPlugin,
	image *images.ImagePlugin,
	tiled *tiled.TiledPlugin) *Onyx {

	o := &Onyx{
		game:       game,
		animations: animations,
		ecs:        ecs,
		images:     image,
		tiled:      tiled,
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
