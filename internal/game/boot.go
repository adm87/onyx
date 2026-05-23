package game

import (
	"context"

	"github.com/adm87/onyx/pkg/engine"
)

func Boot() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	onyx := engine.NewGame(
		engine.WithTitle("Onyx"),
		engine.WithScreenSize(1280, 720),
		engine.WithScreenScale(engine.ScreenScaleFill),
		engine.WithInitialScene(GameSceneIDSplashScreen),
	).WithContext(ctx)

	addAssetAdapters(onyx)
	addScenes(onyx)

	return onyx.Start()
}
