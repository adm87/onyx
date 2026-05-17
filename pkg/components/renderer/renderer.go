package renderer

import (
	"image/color"

	"github.com/adm87/onyx/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RendererData struct {
	Visible bool
	Layer   uint16
}

type ImageData struct {
	Ref   *ebiten.Image
	Color color.RGBA
}

type PathData struct {
	Points []geom.Vec2
	Color  color.RGBA
}

var (
	Renderer = donburi.NewComponentType[RendererData](RendererData{Visible: true})
	Image    = donburi.NewComponentType[ImageData](ImageData{Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}})
	Path     = donburi.NewComponentType[PathData](PathData{Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}})
)

var (
	ImageRenderer = []donburi.IComponentType{
		Renderer,
		Image,
	}
	PathRenderer = []donburi.IComponentType{
		Renderer,
		Path,
	}
)

func GetRenderer(entry *donburi.Entry) *RendererData {
	return Renderer.Get(entry)
}

func GetImage(entry *donburi.Entry) *ImageData {
	return Image.Get(entry)
}

func GetPath(entry *donburi.Entry) *PathData {
	return Path.Get(entry)
}
