package main

import (
	"context"
	"errors"
	"image/color"
	"os"
	"os/signal"
	"syscall"

	"github.com/adm87/onyx/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func gameUpdate(ctx game.Context) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	delta := ctx.Time().Delta32()
	fixed := ctx.Time().FixedDelta32()
	steps := ctx.Time().FixedSteps()

	ctx.Logger().Debug("delta: %f, fixed delta: %f, steps: %d", delta, fixed, steps)
	return nil
}

func gameDraw(ctx game.Context) error {
	ebitenutil.DebugPrint(ctx.Screen().Image(), "Hello, Onyx")
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s := game.NewShell(
		game.WithContext(ctx),
		game.WithTitle("Onyx Game"),
		game.WithSize(1280, 720),
		game.WithOnUpdate(gameUpdate),
		game.WithOnDraw(gameDraw),
		game.WithClearColor(color.RGBA{100, 149, 237, 255}),
		game.WithScreenFilter(ebiten.FilterPixelated),
	)

	if err := s.Run(); err != nil && !errors.Is(err, context.Canceled) {
		s.Context().Logger().Error("game exited with error: %v", err)
		os.Exit(1)
	}
}
