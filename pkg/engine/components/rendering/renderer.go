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

var (
	Filter   = donburi.NewComponentType[ebiten.Filter](ebiten.FilterNearest)
	Renderer = donburi.NewComponentType[RendererData](RendererData{Visible: true})
	Anchor   = donburi.NewComponentType[geom.Vec2]()
	Color    = donburi.NewComponentType[color.RGBA](color.RGBA{R: 255, G: 255, B: 255, A: 255})
)

func GetFilter(entry *donburi.Entry) ebiten.Filter {
	if !entry.HasComponent(Filter) {
		return ebiten.FilterNearest
	}
	return *Filter.Get(entry)
}

func SetFilter(entry *donburi.Entry, filter ebiten.Filter) {
	if !entry.HasComponent(Filter) {
		entry.AddComponent(Filter)
	}
	donburi.SetValue(entry, Filter, filter)
}

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
	renderer := Renderer.Get(entry)
	renderer.Visible = visible
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
	renderer := Renderer.Get(entry)
	renderer.Layer = layer
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
	renderer := Renderer.Get(entry)
	renderer.ZIndex = zIndex
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
