package aseprite

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

const AsepriteRendererType rendering.RendererType = "aseprite_renderer"

type QueryCallback func(*donburi.Entry, geom.Vec2, geom.Vec2, float64)

type AsepriteRendererAdapter struct {
	rendererTypes  []rendering.RendererType
	renderingTasks []engine.RenderingTask
}

func NewAsepriteRendererAdapter() *AsepriteRendererAdapter {
	return &AsepriteRendererAdapter{
		rendererTypes: []rendering.RendererType{
			AsepriteRendererType,
		},
		renderingTasks: make([]engine.RenderingTask, 0, 1),
	}
}

func (a *AsepriteRendererAdapter) SupportedRendererTypes() []rendering.RendererType {
	return a.rendererTypes
}

func (a *AsepriteRendererAdapter) GetRenderingTasks(entry *donburi.Entry, viewport geom.AABB, viewMatrix ebiten.GeoM) []engine.RenderingTask {
	a.renderingTasks = a.renderingTasks[:0]

	return a.renderingTasks
}
