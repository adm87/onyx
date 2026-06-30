package onyx

import (
	"github.com/adm87/onyx/pkg/ecs"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/aseprite"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/adm87/onyx/pkg/plugins/tiled"
	"github.com/yohamta/donburi"
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

	ecs      *ecs.DonburiECS
	aseprite *aseprite.AsepritePlugin
	images   *images.ImagePlugin
	tiled    *tiled.TiledPlugin
}

func NewGame(
	game engine.Game,
	ecs *ecs.DonburiECS,
	animations *aseprite.AsepritePlugin,
	image *images.ImagePlugin,
	tiled *tiled.TiledPlugin) *Onyx {

	o := &Onyx{
		game:     game,
		aseprite: animations,
		ecs:      ecs,
		images:   image,
		tiled:    tiled,
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

func (o *Onyx) AddEntries(entries ...*donburi.Entry) {
	for _, entry := range entries {
		o.ecs.AddEntry(entry)
	}
}

func (o *Onyx) RemoveEntries(entries ...*donburi.Entry) {
	for _, entry := range entries {
		o.ecs.RemoveEntry(entry)
	}
}

func (o *Onyx) UpdateEntries(entries ...*donburi.Entry) {
	for _, entry := range entries {
		o.ecs.UpdateEntry(entry)
	}
}
