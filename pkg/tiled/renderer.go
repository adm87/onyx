package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type TiledRenderingAdapter struct {
	renderingTasks []engine.RenderingTask
	rendererTypes  []rendering.RendererType
}

func NewTiledRenderingAdapter() *TiledRenderingAdapter {
	return &TiledRenderingAdapter{
		renderingTasks: make([]engine.RenderingTask, 0, 10),
		rendererTypes:  []rendering.RendererType{TiledRendererType},
	}
}

func (a *TiledRenderingAdapter) SupportedRendererTypes() []rendering.RendererType {
	return a.rendererTypes
}

func (a *TiledRenderingAdapter) GetRenderingTasks(entry *donburi.Entry, viewport geom.AABB, viewMatrix ebiten.GeoM) []engine.RenderingTask {
	a.renderingTasks = a.renderingTasks[:0]
	return a.renderingTasks
}
