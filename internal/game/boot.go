package game

import (
	"context"
	"path/filepath"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/internal/game/onyx"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/plugins/aseprite"
	"github.com/adm87/onyx/pkg/plugins/ecs"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/adm87/onyx/pkg/plugins/tiled"
	"github.com/hajimehoshi/ebiten/v2"
)

func Boot() error {
	args, err := cli.ParseArgs()
	assert.Fatal(err)

	path, err := filepath.Abs(args.RootDir)
	assert.Fatal(err)

	content.InitContentDirectories(path)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := engine.NewGame(
		engine.WithTitle("Onyx"),
		engine.WithScreenSize(1280, 720),
		engine.WithScreenScale(engine.ScreenScaleFill),
		engine.WithFullscreen(args.Fullscreen),
		engine.WithInitialScene(onyx.SplashScreenSceneID),
		engine.WithFilter(ebiten.FilterNearest),
	).WithContext(ctx)

	assets := g.Assets()
	screen := g.Screen()
	logger := g.Logger()
	renderer := g.Renderer()

	asepritePlugin := aseprite.NewAsepritePlugin(
		logger,
	)

	imagePlugin := images.NewImagePlugin()

	tiledPlugin := tiled.NewTiledPlugin(
		logger,
		imagePlugin.Assets(),
	)

	ecsPlugin := ecs.NewDonburiECSPlugin(
		screen,
		logger,
		imagePlugin.Assets(),
		tiledPlugin.Assets(),
	)

	assets.AddAssetAdapter(imagePlugin.Assets())
	assets.AddAssetAdapter(tiledPlugin.Assets())

	renderer.UsePipeline(ecsPlugin.RenderPipeline())

	return onyx.NewGame(g,
		asepritePlugin,
		ecsPlugin,
		imagePlugin,
		tiledPlugin,
	).Start()
}
