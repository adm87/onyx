package game

import (
	"os"

	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/internal/game/scenes/gameplay"
	"github.com/adm87/onyx/internal/game/scenes/splashscreen"
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

	if err := registerAssetAdapters(onyx); err != nil {
		return err
	}
	if err := registerGameScenes(onyx); err != nil {
		return err
	}
	if err := registerSceneTransitions(onyx); err != nil {
		return err
	}

	args, err := cli.ParseArgs(os.Args[0], os.Args[1:])
	if err != nil {
		onyx.Logger().Error("failed to parse command line arguments: %v", err)
		return err
	}

	ebiten.SetWindowTitle("Onyx " + version)
	ebiten.SetWindowSize(width, height)
	ebiten.SetFullscreen(args.Fullscreen)

	onyx.Screen().SetResizeMode(engine.ScreenResizeByHeight)
	onyx.Scenes().Start(splashscreen.SceneID)

	return onyx.Start()
}

func registerAssetAdapters(onyx engine.Game) error {
	return onyx.Assets().RegisterAdapters(
		images.NewImageAdapter(onyx.Logger()),
	)
}

func registerGameScenes(onyx engine.Game) error {
	return onyx.Scenes().RegisterScenes(
		splashscreen.NewScene(onyx),
		gameplay.NewScene(onyx),
	)
}

func registerSceneTransitions(onyx engine.Game) error {
	if err := onyx.Scenes().RegisterTransitions(
		splashscreen.SceneID,
		engine.SceneTransitions{
			splashscreen.CompleteExitCode: gameplay.SceneID,
		},
	); err != nil {
		return err
	}
	return nil
}
