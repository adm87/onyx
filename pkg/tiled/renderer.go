package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type renderingAdapter struct {
	renderingTasks []engine.RenderingTask
}

func newRenderingAdapter() *renderingAdapter {
	return &renderingAdapter{}
}

func (a *renderingAdapter) GetRenderingTasks(entry *donburi.Entry, viewport geom.AABB, viewMatrix ebiten.GeoM) []engine.RenderingTask {
	a.renderingTasks = a.renderingTasks[:0]

	return a.renderingTasks
}
