package game

import (
	"os"

	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/internal/game/scenes"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	width  = 1280
	height = 720
)

func Boot(version string) error {
	onyx := engine.NewGame(width, height)

	images.Register(onyx.Assets(), onyx.Logger())
	scenes.Register(onyx)

	args, err := cli.ParseArgs(os.Args[0], os.Args[1:])
	if err != nil {
		onyx.Logger().Error("failed to parse command line arguments: %v", err)
		return err
	}

	ebiten.SetWindowTitle("Onyx " + version)
	ebiten.SetWindowSize(width, height)
	ebiten.SetFullscreen(args.Fullscreen)

	onyx.Screen().SetResizeMode(engine.ScreenResizeByHeight)
	onyx.Scenes().Start(scenes.SplashScreenSceneID)

	return onyx.Start()
}
