package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/ecs"
	"github.com/adm87/onyx/pkg/plugins/ecs/renderer"
	"github.com/adm87/onyx/pkg/plugins/ecs/transform"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/yohamta/donburi"
)

var pluginID = engine.TypeHash[TiledPlugin]()

func PluginID() uint64 {
	return pluginID
}

type TiledPlugin interface {
	engine.Plugin

	Assets() *TiledAssets
	CreateTilemap(world donburi.World, opts ...TilemapOption) *donburi.Entry
}

type plugin struct {
	assets   *TiledAssets
	renderer *TiledECSRenderer

	rendererType uint64
}

func NewPlugin() TiledPlugin {
	assets := NewTiledAssets()
	renderer := NewTiledECSRenderer(assets)
	return &plugin{
		assets:   assets,
		renderer: renderer,
	}
}

func (p *plugin) OnRegister(game engine.Game) {
	game.Assets().AddAdapter(p.assets)

	ecsPlugin := engine.GetPlugin[ecs.ECSPlugin](game, ecs.PluginID())
	p.rendererType = ecsPlugin.RenderPipeline().AddAdapter(p.renderer)

	imagePlugin := engine.GetPlugin[images.ImagePlugin](game, images.PluginID())
	p.assets.imageAssets = imagePlugin.Assets()
	p.renderer.imageAssets = imagePlugin.Assets()

	p.renderer.screen = game.Screen()
}

func (p *plugin) ID() uint64 {
	return PluginID()
}

func (p *plugin) Assets() *TiledAssets {
	return p.assets
}

func (p *plugin) CreateTilemap(world donburi.World, opts ...TilemapOption) *donburi.Entry {
	entry := NewTilemap(world, opts...)

	var bounds geom.AABB

	tilemapHandle := GetTilemapHandle(entry)
	if tilemap, exists := p.assets.GetTilemap(tilemapHandle); exists {
		bounds = tilemap.Bounds()
	}

	transform.AddTransform(entry,
		transform.WithBounds(bounds.Min, bounds.Max),
	)

	renderer.AddRenderer(entry,
		renderer.WithRendererType(p.rendererType),
	)

	return entry
}
