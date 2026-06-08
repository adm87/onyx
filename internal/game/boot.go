package game

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/internal/game/scenes"
	"github.com/adm87/onyx/pkg/aseprite"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
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
		engine.WithScreenSize(1280, 720),
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
	if err := images.RegisterPackage(); err != nil {
		return fmt.Errorf("failed to register images package: %w", err)
	}
	if err := tiled.RegisterPackage(); err != nil {
		return fmt.Errorf("failed to register tiled package: %w", err)
	}
	if err := aseprite.RegisterPackage(); err != nil {
		return fmt.Errorf("failed to register aseprite package: %w", err)
	}

	return nil
}
