package engine

import (
	"image/color"

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
	DrawOptions() *ebiten.DrawImageOptions
	SafeArea() *SafeArea
	Clear()
	HandleLayout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	ResizeBuffer(width, height int)
	RestoreBuffer()
	SetBackgroundColor(color color.RGBA)
	SetFilterMode(filter ebiten.Filter)
	SetResizeMode(resizeMode ScreenResizeMode)
}

type screen struct {
	logger Logger
	color  color.RGBA

	buffer   *ebiten.Image
	options  *ebiten.DrawImageOptions
	safeArea *SafeArea

	originalW, originalH int
	logicalW, logicalH   int
	lastW, lastH         int
	isDirty              bool
	scale                float64

	resizeMode ScreenResizeMode
}

func NewScreen(width, height int, logger Logger) Screen {
	if width < 1 || height < 1 {
		panic("screen size must be greater than 0")
	}
	return &screen{
		buffer:    ebiten.NewImage(width, height),
		originalW: width,
		originalH: height,
		safeArea:  &SafeArea{},
		options:   &ebiten.DrawImageOptions{},
	}
}

func (s *screen) Buffer() *ebiten.Image {
	return s.buffer
}

func (s *screen) DrawOptions() *ebiten.DrawImageOptions {
	return s.options
}

func (s *screen) SafeArea() *SafeArea {
	return s.safeArea
}

func (s *screen) Clear() {
	s.buffer.Fill(s.color)
}

func (s *screen) HandleLayout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if outsideWidth != s.lastW || outsideHeight != s.lastH || s.isDirty {
		s.recalculateLayout(outsideWidth, outsideHeight)
	}
	return s.logicalW, s.logicalH
}

func (s *screen) ResizeBuffer(width, height int) {
	if s.buffer.Bounds().Dx() != width || s.buffer.Bounds().Dy() != height {
		s.buffer.Deallocate()
		s.buffer = ebiten.NewImage(width, height)
		s.recalculateLayout(s.lastW, s.lastH)
	}
}

func (s *screen) RestoreBuffer() {
	s.ResizeBuffer(s.originalW, s.originalH)
}

func (s *screen) SetBackgroundColor(color color.RGBA) {
	s.color = color
}

func (s *screen) SetFilterMode(filter ebiten.Filter) {
	s.options.Filter = filter
}

func (s *screen) SetResizeMode(resizeMode ScreenResizeMode) {
	s.resizeMode = resizeMode
}

func (s *screen) recalculateLayout(outsideWidth, outsideHeight int) {
	outsideRatio := float64(outsideWidth) / float64(outsideHeight)
	bufferRatio := float64(s.buffer.Bounds().Dx()) / float64(s.buffer.Bounds().Dy())

	switch s.resizeMode {
	case ScreenResizeByWidth:
		s.resizeByWidth(outsideWidth)
	case ScreenResizeByHeight:
		if outsideRatio > bufferRatio {
			s.resizeByHeight(outsideWidth, outsideHeight)
		} else {
			s.resizeByWidth(outsideWidth)
		}
	case ScreenResizeStretch:
		s.resizeStretch(outsideWidth, outsideHeight)
	}

	s.lastW = outsideWidth
	s.lastH = outsideHeight
	s.isDirty = false
}

func (s *screen) resizeByWidth(outsideWidth int) {
	b := s.buffer.Bounds()

	s.logicalW = b.Dx()
	s.logicalH = b.Dy()

	s.scale = float64(outsideWidth) / float64(b.Dx())
	s.options.GeoM.Reset()

	s.safeArea.minX = 0
	s.safeArea.minY = 0
	s.safeArea.maxX = float64(b.Dx())
	s.safeArea.maxY = float64(b.Dy())
}

func (s *screen) resizeByHeight(outsideWidth, outsideHeight int) {
	b := s.buffer.Bounds()
	bufW, bufH := float64(b.Dx()), float64(b.Dy())

	s.scale = float64(outsideHeight) / bufH

	s.logicalH = int(bufH)
	s.logicalW = int(float64(outsideWidth) / s.scale)

	offX := (float64(s.logicalW) - bufW) / 2

	s.options.GeoM.Reset()
	s.options.GeoM.Translate(offX, 0)

	s.safeArea.minX = -offX
	s.safeArea.minY = 0
	s.safeArea.maxX = offX + bufW
	s.safeArea.maxY = bufH
}

func (s *screen) resizeStretch(outsideWidth, outsideHeight int) {
	b := s.buffer.Bounds()
	bufW, bufH := float64(b.Dx()), float64(b.Dy())

	s.logicalW, s.logicalH = outsideWidth, outsideHeight
	s.scale = float64(outsideWidth) / bufW

	s.options.GeoM.Reset()
	s.options.GeoM.Scale(float64(outsideWidth)/bufW, float64(outsideHeight)/bufH)

	s.safeArea.minX = 0
	s.safeArea.minY = 0
	s.safeArea.maxX = float64(bufW)
	s.safeArea.maxY = float64(bufH)
}
