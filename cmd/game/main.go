package main

import (
	"context"
	"errors"
	"image/color"
	"os"
	"os/signal"
	"syscall"

	"github.com/adm87/onyx/internal/game"
	"github.com/adm87/onyx/pkg/app"
)

var version = "0.0.0-unreleased"

func main() {
	cfg := app.NewConfig(
		app.WithName("Onyx Game"),
		app.WithVersion(version),
		app.WithScreenSize(800, 600),
		app.WithFullscreen(false),
		app.WithLogOutput(os.Stdout),
		app.WithLogLevel(app.LogLevelDebug),
		app.WithScreenClearColor(color.RGBA{R: 100, G: 149, B: 237, A: 255}),
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
