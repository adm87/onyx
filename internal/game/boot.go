package game

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/adm87/onyx/internal/game/input/bindings"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
)

func Boot(cfg *engine.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := engine.NewLogger()
	logger.SetLevel(engine.LogLevelDebug)
	logger.Debug("booting game...")

	assets := engine.NewAssets(logger)
	assets.RegisterAdapter(images.NewEbitenImageAdapter(logger))

	input := engine.NewInput(logger)
	input.Bind(
		bindings.NewGameQuitBinding(),
		bindings.NewFullscreenToggleBinding(),
	)
	input.EnableBinding(bindings.QuitBindingID)
	input.EnableBinding(bindings.FullscreenToggleBindingID)

	screen := engine.NewScreen(
		cfg.Width,
		cfg.Height,
		ebiten.FilterLinear,
		engine.ScreenResizeByHeight,
	)

	onyx := newGame(ctx, cfg, logger, input, assets, screen)
	if err := engine.Run(cfg, onyx); err != nil {
		logger.Error("game loop exited with error", "error", err.Error())
		return err
	}

	logger.Debug("game loop exited, shutting down...")
	return nil
}
