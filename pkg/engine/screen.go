package engine

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

type ScreenScaleMode uint8

const (
	ScreenScaleNone ScreenScaleMode = iota
	ScreenScaleFill
)

type Screen interface {
	ResizeBuffer(width, height int)
	RestoreBuffer()
}

type SafeArea struct {
	Min geom.Vec2
	Max geom.Vec2
}

type screen struct {
	originalWidth, originalHeight int
	outsideWidth, outsideHeight   int
	layoutWidth, layoutHeight     int
	scale                         float64
	isDirty                       bool
	scaleMode                     ScreenScaleMode
	backgroundColor               color.RGBA
	safeArea                      SafeArea

	options *ebiten.DrawImageOptions
	buffer  *ebiten.Image
	logger  *logger
}

func newScreen(
	width, height int,
	scaleMode ScreenScaleMode,
	filter ebiten.Filter,
	backgroundColor color.RGBA,
	logger *logger,
) *screen {

	return &screen{
		originalWidth:   width,
		originalHeight:  height,
		scaleMode:       scaleMode,
		backgroundColor: backgroundColor,
		buffer:          ebiten.NewImage(width, height),
		options: &ebiten.DrawImageOptions{
			Filter: filter,
		},
		safeArea: SafeArea{
			Min: geom.Vec2{X: 0, Y: 0},
			Max: geom.Vec2{X: float64(width), Y: float64(height)},
		},
		logger:  logger,
		isDirty: true,
	}
}

func (s *screen) Layout(outsideWidth, outsideHeight int) (int, int) {
	if s.isDirty || s.outsideWidth != outsideWidth || s.outsideHeight != outsideHeight {
		s.calculateScreenScale(outsideWidth, outsideHeight)
		s.isDirty = false
	}
	return s.layoutWidth, s.layoutHeight
}

func (s *screen) ResizeBuffer(width, height int) {
	bufWidth, bufHeight := s.buffer.Bounds().Dx(), s.buffer.Bounds().Dy()
	if bufWidth == width && bufHeight == height {
		return
	}

	s.logger.Info("Resizing screen buffer from %dx%d to %dx%d", bufWidth, bufHeight, width, height)

	s.buffer.Deallocate()
	s.buffer = ebiten.NewImage(width, height)

	s.isDirty = true
}

func (s *screen) RestoreBuffer() {
	s.ResizeBuffer(s.originalWidth, s.originalHeight)
}

func (s *screen) calculateScreenScale(outsideWidth, outsideHeight int) {
	bufWidth, bufHeight := s.buffer.Bounds().Dx(), s.buffer.Bounds().Dy()

	switch s.scaleMode {
	case ScreenScaleNone:
		s.scale = 1
		s.scaleLayout(bufWidth, bufHeight, outsideWidth, outsideHeight)

	case ScreenScaleFill:
		s.scale = max(
			float64(outsideWidth)/float64(bufWidth),
			float64(outsideHeight)/float64(bufHeight),
		)
		s.scaleLayout(bufWidth, bufHeight, outsideWidth, outsideHeight)
	}

	s.options.GeoM.Reset()
	s.options.GeoM.Translate(-s.safeArea.Min.X, -s.safeArea.Min.Y)

	s.outsideWidth = outsideWidth
	s.outsideHeight = outsideHeight

	s.logger.Debug("Calculated screen scale: %.2f (outside: %dx%d, buffer: %dx%d, layout: %dx%d, safe area: (%.2f, %.2f) - (%.2f, %.2f))",
		s.scale, outsideWidth, outsideHeight, bufWidth, bufHeight, s.layoutWidth, s.layoutHeight,
		s.safeArea.Min.X, s.safeArea.Min.Y, s.safeArea.Max.X, s.safeArea.Max.Y)
}

func (s *screen) scaleLayout(bufWidth, bufHeight, outsideWidth, outsideHeight int) {
	s.layoutWidth = int(float64(outsideWidth) / s.scale)
	s.layoutHeight = int(float64(outsideHeight) / s.scale)
	s.safeArea.Min.X = (float64(bufWidth) - float64(s.layoutWidth)) / 2
	s.safeArea.Min.Y = (float64(bufHeight) - float64(s.layoutHeight)) / 2
	s.safeArea.Max.X = s.safeArea.Min.X + float64(s.layoutWidth)
	s.safeArea.Max.Y = s.safeArea.Min.Y + float64(s.layoutHeight)
}
