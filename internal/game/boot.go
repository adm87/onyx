package game

import (
	"image/color"
	"os"

	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func Boot(version string) error {
	const (
		width  = 800
		height = 600
	)

	world := donburi.NewWorld()

	logger := engine.NewLogger(os.Stdout)
	assets := engine.NewAssets(logger)
	scenes := engine.NewScenes(world, logger)
	time := engine.NewTime()

	screen := engine.NewScreen(width, height, logger)
	screen.SetBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255})

	args, err := cli.ParseArgs(os.Args[0], os.Args[1:])
	if err != nil {
		logger.Error("failed to parse arguments: %v", err)
		return err
	}

	logger.Debug("setting window")
	setupWindow("Onyx", version, width, height, args.Fullscreen)

	logger.Debug("running game")
	return ebiten.RunGame(&game{
		assets: assets,
		logger: logger,
		scenes: scenes,
		screen: screen,
		time:   time,
	})
}

func setupWindow(name, version string, width, height int, fullsreen bool) {
	ebiten.SetWindowTitle(name + " " + version)
	ebiten.SetWindowSize(width, height)
	ebiten.SetFullscreen(fullsreen)
}
