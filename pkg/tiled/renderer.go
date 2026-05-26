package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
)

type TiledRenderingAdapter struct {
	logger         engine.Logger
	renderingTasks []engine.RenderTask
}

func NewTiledRenderingAdapter(logger engine.Logger) *TiledRenderingAdapter {
	return &TiledRenderingAdapter{
		logger:         logger,
		renderingTasks: make([]engine.RenderTask, 0, 100),
	}
}

func (a *TiledRenderingAdapter) GetRenderTasks(world donburi.World) []engine.RenderTask {
	a.renderingTasks = a.renderingTasks[:0]

	return a.renderingTasks
}
