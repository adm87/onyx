package game

import (
	"context"
	"path/filepath"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/pkg/engine"
)

func Boot() error {
	args, err := cli.ParseArgs()
	if err != nil {
		return err
	}

	path, err := filepath.Abs(args.RootDir)
	if err != nil {
		return err
	}

	content.InitContentDirectories(path)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	onyx := engine.NewGame(
		engine.WithTitle("Onyx"),
		engine.WithScreenSize(1280, 720),
		engine.WithScreenScale(engine.ScreenScaleFill),
		engine.WithInitialScene(SplashScreenSceneID),
	).WithContext(ctx)

	addAssetAdapters(onyx)
	addRenderingSystems(onyx)
	addScenes(onyx)

	if err := content.LoadDefaultContent(onyx.Assets(), onyx.Logger()); err != nil {
		return err
	}

	return onyx.Start()
}
