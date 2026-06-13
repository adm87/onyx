package game

import (
	"context"
	"image/color"
	"path/filepath"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/internal/game/onyx"
	"github.com/adm87/onyx/pkg/aseprite"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
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
		engine.WithBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255}),
		engine.WithFullscreen(args.Fullscreen),
		engine.WithInitialScene(onyx.GameplaySceneID),
		engine.WithFilter(ebiten.FilterNearest),
	).WithContext(ctx)

	assets := g.Assets()
	renderer := g.Renderer()
	screen := g.Screen()

	imageModule := images.NewModule(
		assets,
		renderer,
	)
	tiledModule := tiled.NewModule(
		assets,
		renderer,
		screen,
		imageModule,
	)
	asepriteModule := aseprite.NewAsepriteModule(
		imageModule,
	)

	return onyx.NewGame(g,
		imageModule,
		tiledModule,
		asepriteModule,
	).Start()
}
