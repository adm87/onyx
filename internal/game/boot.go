package game

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/adm87/onyx-game/content"
	"github.com/adm87/onyx-game/internal/game/cli"
	"github.com/adm87/onyx-game/internal/game/scenes"
	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/adm87/onyx-game/pkg/images"
	"github.com/adm87/onyx-game/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
)

func Boot() error {
	args, err := cli.ParseArgs()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	onyx := engine.NewGame(
		engine.WithTitle("Onyx"),
		engine.WithScreenSize(800, 600),
		engine.WithFullscreen(args.Fullscreen),
		engine.WithScreenScale(engine.ScreenScaleFill),
		engine.WithInitialScene(scenes.GameplaySceneID),
		engine.WithFilter(ebiten.FilterNearest),
	).WithContext(ctx)

	if err := registerPackages(onyx); err != nil {
		return err
	}

	path, err := filepath.Abs(args.RootDir)
	if err != nil {
		return err
	}

	content.InitContentDirectories(path)

	if err := content.LoadDefaultContent(onyx.Assets(), onyx.Logger()); err != nil {
		return err
	}

	scenes.AddScenes(onyx)
	return onyx.Start()
}

func registerPackages(onyx engine.Game) error {
	assets := onyx.Assets()
	screen := onyx.Screen()
	renderer := onyx.Renderer()
	logger := onyx.Logger()

	// NOTE: The order of package registration matters, as some packages may depend on others being registered first.

	if err := images.RegisterPackage(assets, renderer); err != nil {
		return fmt.Errorf("failed to register images package: %w", err)
	}
	if err := tiled.RegisterPackage(assets, renderer, screen, logger); err != nil {
		return fmt.Errorf("failed to register tiled package: %w", err)
	}

	return nil
}
