package engine

import "github.com/hajimehoshi/ebiten/v2"

type ScreenResizeMode uint8

const (
	ScreenResizeByWidth ScreenResizeMode = iota
	ScreenResizeByHeight
	ScreenResizeStretch
)

type SafeArea struct {
	minX, minY float64
	maxX, maxY float64
}

func (s SafeArea) Min() (float64, float64) {
	return s.minX, s.minY
}

func (s SafeArea) Max() (float64, float64) {
	return s.maxX, s.maxY
}

type Screen struct {
	img      *ebiten.Image
	opts     *ebiten.DrawImageOptions
	safeArea *SafeArea

	logicalW, logicalH int
	lastW, lastH       int
	isDirty            bool
	scale              float64

	resizeMode ScreenResizeMode
}

func NewScreen(width, height int, filter ebiten.Filter, resizeMode ScreenResizeMode) *Screen {
	return &Screen{
		img: ebiten.NewImage(width, height),
		opts: &ebiten.DrawImageOptions{
			Filter: filter,
		},
		resizeMode: resizeMode,
		safeArea:   &SafeArea{},
	}
}

func (s *Screen) Buffer() *ebiten.Image {
	return s.img
}

func (s *Screen) Options() *ebiten.DrawImageOptions {
	return s.opts
}

func (s *Screen) ResizeBuffer(width, height int) {
	if s.img.Bounds().Dx() != width || s.img.Bounds().Dy() != height {
		s.img.Deallocate()
		s.img = ebiten.NewImage(width, height)
		s.recalculateLayout(s.lastW, s.lastH)
	}
}

func (s *Screen) Scale() float64 {
	return s.scale
}

func (s *Screen) SafeArea() *SafeArea {
	return s.safeArea
}

func (s *Screen) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth == s.logicalW && outsideHeight == s.logicalH && !s.isDirty {
		return s.logicalW, s.logicalH
	}

	s.recalculateLayout(outsideWidth, outsideHeight)
	return s.logicalW, s.logicalH
}

func (s *Screen) recalculateLayout(outsideWidth, outsideHeight int) {
	outsideRatio := float64(outsideWidth) / float64(outsideHeight)
	bufferRatio := float64(s.img.Bounds().Dx()) / float64(s.img.Bounds().Dy())

	if s.resizeMode == ScreenResizeByWidth {
		s.resizeByWidth(outsideWidth, outsideHeight)
	} else if s.resizeMode == ScreenResizeByHeight {
		if outsideRatio > bufferRatio {
			s.resizeByWidth(outsideWidth, outsideHeight)
		} else {
			s.resizeByHeight(outsideWidth, outsideHeight)
		}
	} else if s.resizeMode == ScreenResizeStretch {
		s.resizeStretch(outsideWidth, outsideHeight)
	}

	s.lastW = outsideWidth
	s.lastH = outsideHeight
	s.isDirty = false
}

func (s *Screen) resizeByWidth(outsideWidth, outsideHeight int) {
	b := s.img.Bounds()

	s.logicalW = b.Dx()
	s.logicalH = b.Dy()

	s.scale = float64(outsideWidth) / float64(b.Dx())

	s.opts.GeoM.Reset()

	s.safeArea.minX = 0
	s.safeArea.minY = 0
	s.safeArea.maxX = float64(b.Dx())
	s.safeArea.maxY = float64(b.Dy())
}

func (s *Screen) resizeByHeight(outsideWidth, outsideHeight int) {
	b := s.img.Bounds()
	bufW, bufH := float64(b.Dx()), float64(b.Dy())

	s.scale = float64(outsideHeight) / bufH

	s.logicalH = int(bufH)
	s.logicalW = int(float64(outsideWidth) / s.scale)

	offX := (float64(s.logicalW) - bufW) / 2

	s.opts.GeoM.Reset()
	s.opts.GeoM.Translate(offX, 0)

	s.safeArea.minX = -offX
	s.safeArea.minY = 0
	s.safeArea.maxX = offX + bufW
	s.safeArea.maxY = bufH
}

func (s *Screen) resizeStretch(outsideWidth, outsideHeight int) {
	b := s.img.Bounds()
	bufW, bufH := float64(b.Dx()), float64(b.Dy())

	s.logicalW, s.logicalH = outsideWidth, outsideHeight
	s.scale = float64(outsideWidth) / bufW

	s.opts.GeoM.Reset()
	s.opts.GeoM.Scale(float64(outsideWidth)/bufW, float64(outsideHeight)/bufH)

	s.safeArea.minX = 0
	s.safeArea.minY = 0
	s.safeArea.maxX = float64(bufW)
	s.safeArea.maxY = float64(bufH)
}
