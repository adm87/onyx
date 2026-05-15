package game

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

type shell struct {
	logger engine.Logger
}

func (s *shell) Update() error {
	return nil
}

func (s *shell) Draw(screen *ebiten.Image) {

}

func (s *shell) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
