package tiled

import (
	"github.com/adm87/onyx/pkg/ecs/renderer"
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/yohamta/donburi"
)

type TiledPlugin struct {
	assets   *TiledAssets
	renderer *TiledECSRenderer
}

func NewTiledPlugin(screen engine.Screen, images *images.ImageAssets) *TiledPlugin {
	assets := NewTiledAssets(images)
	return &TiledPlugin{
		assets:   assets,
		renderer: NewTiledECSRenderer(screen, images, assets),
	}
}

func (t *TiledPlugin) Assets() *TiledAssets {
	return t.assets
}

func (t *TiledPlugin) Renderer() *TiledECSRenderer {
	return t.renderer
}

func (t *TiledPlugin) CreateTilemap(world donburi.World, opts ...TilemapOption) *donburi.Entry {
	entry := NewTilemap(world, opts...)

	var bounds geom.AABB

	tilemapHandle := GetTilemapHandle(entry)
	if tilemap, exists := t.assets.GetTilemap(tilemapHandle); exists {
		bounds = tilemap.Bounds()
	}

	transform.AddTransform(entry,
		transform.WithBounds(bounds.Min, bounds.Max),
	)

	renderer.AddRenderer(entry,
		renderer.WithRendererType(t.renderer.adapterIndex),
	)

	return entry
}
