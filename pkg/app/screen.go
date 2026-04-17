package app

import "github.com/hajimehoshi/ebiten/v2"

type Screen struct {
	buffer *ebiten.Image
}

func NewScreen(width, height int) *Screen {
	return &Screen{
		buffer: ebiten.NewImage(width, height),
	}
}

func (s *Screen) layout(outsideWidth, outsideHeight int) (innerWidth, innerHeight int) {
	return s.buffer.Bounds().Dx(), s.buffer.Bounds().Dy()
}
