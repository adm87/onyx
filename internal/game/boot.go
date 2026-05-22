package game

import (
	"context"
	"image/color"

	"github.com/adm87/onyx/pkg/engine"
)

func Boot() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	game := engine.NewGame(
		engine.WithTitle("Onyx"),
		engine.WithScreenSize(1280, 720),
		engine.WithScreenScale(engine.ScreenScaleFill),
		engine.WithBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255}),
	).WithContext(ctx)

	return game.Start()
}
