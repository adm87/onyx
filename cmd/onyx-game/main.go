package main

import (
	"context"
	"errors"
	"image/color"
	"os"
	"os/signal"
	"syscall"

	onyx "github.com/adm87/onyx/internal/game"
	"github.com/adm87/onyx/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
)

var version = "0.0.0-unreleased" // This will be set at build time using -ldflags="-X main.version=1.0.0"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s := game.NewShell(
		game.WithContext(ctx),
		game.WithTitle("Onyx Game"),
		game.WithVersion(version),
		game.WithSize(1280, 720),
		game.WithClearColor(color.RGBA{100, 149, 237, 255}),
		game.WithScreenFilter(ebiten.FilterPixelated),
		game.WithModel(onyx.NewModel()),
		game.WithOnStart(onyx.OnStart),
		game.WithOnUpdate(onyx.OnUpdate),
		game.WithOnFixedUpdate(onyx.OnFixedUpdate),
		game.WithOnLateUpdate(onyx.OnLateUpdate),
		game.WithOnDraw(onyx.OnDraw),
	)

	if err := s.Run(); err != nil && !errors.Is(err, context.Canceled) {
		s.Context().Logger().Error("game exited with error: %v", err)
		os.Exit(1)
	}
}
