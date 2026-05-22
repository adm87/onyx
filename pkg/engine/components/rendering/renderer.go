package rendering

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RendererData struct {
	Visible bool
	Layer   int
	ZIndex  int
}

type ImageData struct {
	Ref *ebiten.Image
}

var (
	Renderer = donburi.NewComponentType[RendererData](RendererData{Visible: true})
	Image    = donburi.NewComponentType[ImageData]()
	Anchor   = donburi.NewComponentType[geom.Vec2]()
	Color    = donburi.NewComponentType[color.RGBA](color.RGBA{R: 255, G: 255, B: 255, A: 255})
)

func IsVisible(entry *donburi.Entry) bool {
	if !entry.HasComponent(Renderer) {
		return false
	}
	return Renderer.Get(entry).Visible
}

func SetVisible(entry *donburi.Entry, visible bool) {
	if !entry.HasComponent(Renderer) {
		entry.AddComponent(Renderer)
	}
	donburi.SetValue(entry, Renderer, RendererData{
		Visible: visible,
		Layer:   Renderer.Get(entry).Layer,
		ZIndex:  Renderer.Get(entry).ZIndex,
	})
}

func GetLayer(entry *donburi.Entry) int {
	if !entry.HasComponent(Renderer) {
		return 0
	}
	return Renderer.Get(entry).Layer
}

func SetLayer(entry *donburi.Entry, layer int) {
	if !entry.HasComponent(Renderer) {
		entry.AddComponent(Renderer)
	}
	donburi.SetValue(entry, Renderer, RendererData{
		Visible: Renderer.Get(entry).Visible,
		Layer:   layer,
		ZIndex:  Renderer.Get(entry).ZIndex,
	})
}

func GetZIndex(entry *donburi.Entry) int {
	if !entry.HasComponent(Renderer) {
		return 0
	}
	return Renderer.Get(entry).ZIndex
}

func SetZIndex(entry *donburi.Entry, zIndex int) {
	if !entry.HasComponent(Renderer) {
		entry.AddComponent(Renderer)
	}
	donburi.SetValue(entry, Renderer, RendererData{
		Visible: Renderer.Get(entry).Visible,
		Layer:   Renderer.Get(entry).Layer,
		ZIndex:  zIndex,
	})
}

func GetImage(entry *donburi.Entry) *ebiten.Image {
	if !entry.HasComponent(Image) {
		return nil
	}
	return Image.Get(entry).Ref
}

func SetImage(entry *donburi.Entry, img *ebiten.Image) {
	if !entry.HasComponent(Image) {
		entry.AddComponent(Image)
	}
	donburi.SetValue(entry, Image, ImageData{
		Ref: img,
	})
}

func GetAnchor(entry *donburi.Entry) geom.Vec2 {
	if !entry.HasComponent(Anchor) {
		return geom.Vec2{X: 0, Y: 0}
	}
	return *Anchor.Get(entry)
}

func SetAnchor(entry *donburi.Entry, anchor geom.Vec2) {
	if !entry.HasComponent(Anchor) {
		entry.AddComponent(Anchor)
	}
	donburi.SetValue(entry, Anchor, anchor)
}

func GetColor(entry *donburi.Entry) color.RGBA {
	if !entry.HasComponent(Color) {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	return *Color.Get(entry)
}

func SetColor(entry *donburi.Entry, col color.RGBA) {
	if !entry.HasComponent(Color) {
		entry.AddComponent(Color)
	}
	donburi.SetValue(entry, Color, col)
}
