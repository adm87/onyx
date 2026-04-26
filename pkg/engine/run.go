package engine

import "github.com/hajimehoshi/ebiten/v2"

func Run(cfg *Config, game Game) error {
	initializeWindow(cfg)
	return ebiten.RunGame(game)
}

func initializeWindow(cfg *Config) {
	ebiten.SetWindowTitle(cfg.Title)
	ebiten.SetWindowSize(cfg.Width, cfg.Height)
	ebiten.SetFullscreen(cfg.Fullscreen)
}
