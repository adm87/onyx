package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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

type Screen interface {
	Buffer() *ebiten.Image
	Options() *ebiten.DrawImageOptions

	ResizeBuffer(width, height int)
	RestoreBuffer()

	Scale() float64
	SafeArea() *SafeArea

	Min() (float64, float64)
	Max() (float64, float64)

	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)

	SetResizeMode(ScreenResizeMode)
}

type screen struct {
	logger Logger

	img      *ebiten.Image
	opts     *ebiten.DrawImageOptions
	safeArea *SafeArea

	originalW, originalH int
	logicalW, logicalH   int
	lastW, lastH         int
	isDirty              bool
	scale                float64

	resizeMode ScreenResizeMode
}

func NewScreen(width, height int, logger Logger) Screen {
	screen := &screen{
		logger:    logger,
		img:       ebiten.NewImage(width, height),
		opts:      &ebiten.DrawImageOptions{},
		originalW: width,
		originalH: height,
		lastW:     width,
		lastH:     height,
		safeArea:  &SafeArea{},
		isDirty:   true,
	}
	return screen
}

func (s *screen) SetResizeMode(mode ScreenResizeMode) {
	s.resizeMode = mode
	s.isDirty = true
}

func (s *screen) Buffer() *ebiten.Image {
	return s.img
}

func (s *screen) Options() *ebiten.DrawImageOptions {
	return s.opts
}

func (s *screen) ResizeBuffer(width, height int) {
	if s.img.Bounds().Dx() != width || s.img.Bounds().Dy() != height {
		s.img.Deallocate()
		s.img = ebiten.NewImage(width, height)
		s.recalculateLayout(s.lastW, s.lastH)
	}
}

func (s *screen) RestoreBuffer() {
	s.ResizeBuffer(s.originalW, s.originalH)
}

func (s *screen) Scale() float64 {
	return s.scale
}

func (s *screen) SafeArea() *SafeArea {
	return s.safeArea
}

func (s *screen) Min() (float64, float64) {
	return 0, 0
}

func (s *screen) Max() (float64, float64) {
	return float64(s.img.Bounds().Dx()), float64(s.img.Bounds().Dy())
}

func (s *screen) Layout(outsideWidth, outsideHeight int) (int, int) {
	if s.isDirty || s.lastW != outsideWidth || s.lastH != outsideHeight {
		s.recalculateLayout(outsideWidth, outsideHeight)
	}
	return s.logicalW, s.logicalH
}

func (s *screen) recalculateLayout(outsideWidth, outsideHeight int) {
	s.logger.Debug("recalculating layout for %dx%d", outsideWidth, outsideHeight)

	switch s.resizeMode {
	case ScreenResizeByWidth:
		s.logger.Debug("Resizing screen by width")
		s.resizeByWidth(outsideWidth)
	case ScreenResizeByHeight:
		s.logger.Debug("Resizing screen by height")

		outsideRatio := float64(outsideWidth) / float64(outsideHeight)
		bufferRatio := float64(s.img.Bounds().Dx()) / float64(s.img.Bounds().Dy())

		if outsideRatio >= bufferRatio {
			s.resizeByWidth(outsideWidth)
		} else {
			s.resizeByHeight(outsideWidth, outsideHeight)
		}
	case ScreenResizeStretch:
		s.logger.Debug("Resizing screen by stretching")
		s.resizeStretch(outsideWidth, outsideHeight)
	}

	s.lastW = outsideWidth
	s.lastH = outsideHeight
	s.isDirty = false
}

func (s *screen) resizeByWidth(outsideWidth int) {
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

func (s *screen) resizeByHeight(outsideWidth, outsideHeight int) {
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

func (s *screen) resizeStretch(outsideWidth, outsideHeight int) {
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
