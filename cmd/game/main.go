package main

import (
	"context"
	"errors"
	"image/color"
	"os"
	"os/signal"
	"syscall"

	"github.com/adm87/onyx/internal/game"
	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/pkg/app"
	"github.com/hajimehoshi/ebiten/v2"
)

var version = "0.0.0-unreleased"

func main() {
	args := cli.NewArgs()
	if err := args.Parse(os.Args[0], os.Args[1:]); err != nil {
		println("error parsing command-line arguments:", err.Error())
		os.Exit(1)
	}

	cfg := app.NewConfig(
		app.WithName("Onyx Game"),
		app.WithVersion(version),
		app.WithScreenSize(1280, 720),
		app.WithFullscreen(args.Fullscreen),
		app.WithLogOutput(os.Stdout),
		app.WithLogLevel(app.LogLevelDebug),
		app.WithScreenClearColor(color.RGBA{R: 100, G: 149, B: 237, A: 255}),
		app.WithScreenFilter(ebiten.FilterPixelated),
		app.WithScreenResizeMode(app.ScreenResizeStretch),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx, game.New(), cfg); err != nil {
		println("error running application:", err.Error())
		if !errors.Is(err, context.Canceled) {
			os.Exit(1)
		}
	}
}
