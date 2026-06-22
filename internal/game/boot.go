package game

import (
	"context"
	"path/filepath"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/internal/game/onyx"
	"github.com/adm87/onyx/pkg/ecs"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/plugins/aseprite"
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
		engine.WithInitialScene(onyx.GameplaySceneID),
		engine.WithFilter(ebiten.FilterNearest),
	).WithContext(ctx)

	assets := g.Assets()
	screen := g.Screen()
	logger := g.Logger()
	renderer := g.Renderer()

	donburiECS := ecs.NewDonburiECS(
		screen,
		logger,
	)

	imagePlugin := images.NewImagePlugin()

	asepritePlugin := aseprite.NewAsepritePlugin(
		imagePlugin,
	)

	tiledPlugin := tiled.NewTiledPlugin(
		screen,
		imagePlugin.Assets(),
	)

	donburiECS.RenderPipeline().AddAdapters(
		imagePlugin.Renderer(),
		tiledPlugin.Renderer(),
	)

	assets.AddAdapters(imagePlugin.Assets())
	assets.AddAdapters(tiledPlugin.Assets())

	renderer.UsePipeline(donburiECS.RenderPipeline())

	return onyx.NewGame(g, donburiECS,
		asepritePlugin,
		imagePlugin,
		tiledPlugin,
	).Start()
}
