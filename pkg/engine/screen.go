package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Screen interface {
	Buffer() *ebiten.Image
	DrawOptions() *ebiten.DrawImageOptions
	Clear()
	HandleLayout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	ResizeBuffer(width, height int)
	RestoreBuffer()
	SetBackgroundColor(color color.RGBA)
	SetFilterMode(filter ebiten.Filter)
}

type screen struct {
	width, height int
	logger        Logger
	color         color.RGBA

	buffer  *ebiten.Image
	options *ebiten.DrawImageOptions
}

func NewScreen(width, height int, logger Logger) Screen {
	if width < 1 || height < 1 {
		panic("screen size must be greater than 0")
	}

	img := ebiten.NewImage(width, height)
	return &screen{
		width:  width,
		height: height,
		logger: logger,
		buffer: img,
	}
}

func (s *screen) Buffer() *ebiten.Image {
	return s.buffer
}

func (s *screen) DrawOptions() *ebiten.DrawImageOptions {
	return s.options
}

func (s *screen) Clear() {
	s.buffer.Fill(s.color)
}

func (s *screen) HandleLayout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.width, s.height
}

func (s *screen) ResizeBuffer(width, height int) {
	if width < 1 || height < 1 {
		panic("buffer size must be greater than 0")
	}

	if width == s.width && height == s.height {
		return
	}

	s.logger.Debug("resizing screen buffer to %dx%d", width, height)

	s.buffer.Deallocate()
	s.buffer = ebiten.NewImage(width, height)
}

func (s *screen) RestoreBuffer() {
	width, height := s.buffer.Bounds().Dx(), s.buffer.Bounds().Dy()

	if width == s.width && height == s.height {
		return
	}

	s.logger.Debug("restoring screen buffer to %dx%d", s.width, s.height)

	s.buffer.Deallocate()
	s.buffer = ebiten.NewImage(s.width, s.height)
}

func (s *screen) SetBackgroundColor(color color.RGBA) {
	s.color = color
}

func (s *screen) SetFilterMode(filter ebiten.Filter) {
	s.options.Filter = filter
}
