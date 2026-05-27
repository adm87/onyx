package tiled

import (
	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/adm87/onyx-game/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type TiledRenderingAdapter struct {
	tiledAssetAdapter *TiledAssetAdapter
	imageAssetAdapter *images.ImageAssetAdapter

	screen         engine.Screen
	renderingTasks []engine.RenderTask
}

func NewTiledRenderingAdapter(tiledAssetAdapter *TiledAssetAdapter, imageAssetAdapter *images.ImageAssetAdapter, screen engine.Screen) *TiledRenderingAdapter {
	return &TiledRenderingAdapter{
		tiledAssetAdapter: tiledAssetAdapter,
		imageAssetAdapter: imageAssetAdapter,
		screen:            screen,
		renderingTasks:    make([]engine.RenderTask, 0, 10),
	}
}

func (a *TiledRenderingAdapter) GetRenderTasks(world donburi.World, viewMatrix ebiten.GeoM) []engine.RenderTask {
	a.renderingTasks = a.renderingTasks[:0]

	screenMinX, screenMinY := a.screen.SafeArea().Min.XY()
	screenMaxX, screenMaxY := a.screen.SafeArea().Max.XY()

	worldMinX, worldMinY := viewMatrix.Apply(screenMinX, screenMinY)
	worldMaxX, worldMaxY := viewMatrix.Apply(screenMaxX, screenMaxY)

	_ = worldMinX
	_ = worldMinY
	_ = worldMaxX
	_ = worldMaxY

	TiledQuery.Each(world, func(e *donburi.Entry) {

	})

	return a.renderingTasks
}
