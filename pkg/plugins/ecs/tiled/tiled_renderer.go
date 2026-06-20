package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/ecs/renderer"
	imageplugin "github.com/adm87/onyx/pkg/plugins/images"
	tiledplugin "github.com/adm87/onyx/pkg/plugins/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type TiledRenderer struct {
	imageAssets *imageplugin.ImageAssets
	tiledAssets *tiledplugin.TiledAssets
}

func NewTiledRenderer(
	imageAssets *imageplugin.ImageAssets,
	tiledAssets *tiledplugin.TiledAssets) *TiledRenderer {
	return &TiledRenderer{
		imageAssets: imageAssets,
		tiledAssets: tiledAssets,
	}
}

func (r *TiledRenderer) PrepareRenderingTasks(
	entry *donburi.Entry,
	renderer *renderer.RendererModel,
	pool *engine.RenderingPool,
	viewport geom.AABB,
	viewMatrix ebiten.GeoM) []*engine.RenderingTask {

	return nil
}
