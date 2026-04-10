package game

import "github.com/hajimehoshi/ebiten/v2"

type Screen interface {
	Image() *ebiten.Image
	DrawImage(img *ebiten.Image, opts *ebiten.DrawImageOptions)
}

type screenImpl struct {
	img *ebiten.Image
	opt *ebiten.DrawImageOptions
}

func (s *screenImpl) layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.img.Bounds().Dx(), s.img.Bounds().Dy()
}

func (s *screenImpl) Image() *ebiten.Image {
	return s.img
}

func (s *screenImpl) DrawImage(img *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.img.DrawImage(img, opts)
}
