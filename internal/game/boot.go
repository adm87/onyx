package game

import (
	"os"

	"github.com/adm87/onyx/internal/game/cli"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

func Boot(version string) error {
	logger := engine.NewLogger(os.Stdout)

	args, err := cli.ParseArgs(os.Args[0], os.Args[1:])
	if err != nil {
		logger.Error("failed to parse arguments: %v", err)
		return err
	}

	logger.Debug("setting window")
	setupWindow("Onyx", version, 800, 600, args.Fullscreen)

	logger.Debug("running game")
	return ebiten.RunGame(&shell{
		logger: logger,
	})
}

func setupWindow(name, version string, width, height int, fullsreen bool) {
	ebiten.SetWindowTitle(name + " " + version)
	ebiten.SetWindowSize(width, height)
	ebiten.SetFullscreen(fullsreen)
}
