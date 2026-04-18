package app

import "github.com/hajimehoshi/ebiten/v2"

type ScreenResizeMode uint8

const (
	ScreenResizeByWidth ScreenResizeMode = iota
	ScreenResizeByHeight
	ScreenResizeStretch
)

type Screen struct {
	buffer     *ebiten.Image
	opt        *ebiten.DrawImageOptions
	resizeMode ScreenResizeMode

	logicalW, logicalH int
	lastW, lastH       int
	isDirty            bool

	scale   float64
	safeMin [2]float64
	safeMax [2]float64
}

func NewScreen(width, height int, filter ebiten.Filter, resizeMode ScreenResizeMode) *Screen {
	return &Screen{
		buffer: ebiten.NewImage(width, height),
		opt: &ebiten.DrawImageOptions{
			Filter: filter,
		},
		resizeMode: resizeMode,
	}
}

func (s *Screen) Buffer() *ebiten.Image {
	return s.buffer
}

func (s *Screen) ResizeBuffer(width, height int) {
	if s.buffer.Bounds().Dx() != width || s.buffer.Bounds().Dy() != height {
		s.buffer.Deallocate()
		s.buffer = ebiten.NewImage(width, height)
		s.recalculateLayout(s.lastW, s.lastH)
	}
}

func (s *Screen) Scale() float64 {
	return s.scale
}

func (s *Screen) SafeMin() (float64, float64) {
	return s.safeMin[0], s.safeMin[1]
}

func (s *Screen) SafeMax() (float64, float64) {
	return s.safeMax[0], s.safeMax[1]
}

func (s *Screen) layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth == s.logicalW && outsideHeight == s.logicalH && !s.isDirty {
		return s.logicalW, s.logicalH
	}

	s.recalculateLayout(outsideWidth, outsideHeight)
	return s.logicalW, s.logicalH
}

func (s *Screen) recalculateLayout(outsideWidth, outsideHeight int) {
	outsideRatio := float64(outsideWidth) / float64(outsideHeight)
	bufferRatio := float64(s.buffer.Bounds().Dx()) / float64(s.buffer.Bounds().Dy())

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
	b := s.buffer.Bounds()

	s.logicalW = b.Dx()
	s.logicalH = b.Dy()

	s.scale = float64(outsideWidth) / float64(b.Dx())

	s.opt.GeoM.Reset()

	s.safeMin = [2]float64{0, 0}
	s.safeMax = [2]float64{float64(b.Dx()), float64(b.Dy())}
}

func (s *Screen) resizeByHeight(outsideWidth, outsideHeight int) {
	b := s.buffer.Bounds()
	bufW, bufH := float64(b.Dx()), float64(b.Dy())

	s.scale = float64(outsideHeight) / bufH

	s.logicalH = int(bufH)
	s.logicalW = int(float64(outsideWidth) / s.scale)

	offX := (float64(s.logicalW) - bufW) / 2

	s.opt.GeoM.Reset()
	s.opt.GeoM.Translate(offX, 0)

	s.safeMin = [2]float64{-offX, 0}
	s.safeMax = [2]float64{offX + bufW, bufH}
}

func (s *Screen) resizeStretch(outsideWidth, outsideHeight int) {
	b := s.buffer.Bounds()
	bufW, bufH := float64(b.Dx()), float64(b.Dy())

	s.logicalW, s.logicalH = outsideWidth, outsideHeight
	s.scale = float64(outsideWidth) / bufW

	s.opt.GeoM.Reset()
	s.opt.GeoM.Scale(float64(outsideWidth)/bufW, float64(outsideHeight)/bufH)

	s.safeMin = [2]float64{0, 0}
	s.safeMax = [2]float64{float64(bufW), float64(bufH)}
}
