package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RendererData struct {
	Visible bool
	Layer   int
}

type ImageData struct {
	Ref   *ebiten.Image
	Color color.RGBA
}

type PathData struct {
	Points []float64
	Color  color.RGBA
}

var (
	Renderer      = donburi.NewComponentType[RendererData](RendererData{Visible: true})
	ImageRenderer = donburi.NewComponentType[ImageData](ImageData{Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}})
	PathRenderer  = donburi.NewComponentType[PathData](PathData{Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}})
)
