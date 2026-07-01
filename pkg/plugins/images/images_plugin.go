package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/ecs"
	"github.com/adm87/onyx/pkg/plugins/ecs/renderer"
	"github.com/adm87/onyx/pkg/plugins/ecs/transform"
	"github.com/yohamta/donburi"
)

var pluginID = engine.TypeHash[ImagePlugin]()

func PluginID() uint64 {
	return pluginID
}

type ImagePlugin interface {
	engine.Plugin

	Assets() *ImageAssets
	CreateImage(world donburi.World, opts ...Option) *donburi.Entry
}

type plugin struct {
	assets   *ImageAssets
	renderer *ImageECSRenderer

	rendererType uint64
}

func NewPlugin() ImagePlugin {
	assets := NewImageAssets()
	renderer := NewImageECSRenderer(assets)
	return &plugin{
		assets:   assets,
		renderer: renderer,
	}
}

func (p *plugin) OnRegister(game engine.Game) {
	game.Assets().AddAdapter(p.assets)

	ecsPlugin := engine.GetPlugin[ecs.ECSPlugin](game, ecs.PluginID())
	p.rendererType = ecsPlugin.RenderPipeline().AddAdapter(p.renderer)
}

func (p *plugin) ID() uint64 {
	return PluginID()
}

func (p *plugin) Assets() *ImageAssets {
	return p.assets
}

func (p *plugin) CreateImage(world donburi.World, opts ...Option) *donburi.Entry {
	entry := NewImage(world, opts...)

	var bounds geom.AABB

	imgHandle := GetHandle(entry)
	frameIdx := GetFrame(entry)

	if img, exists := p.assets.GetFrame(imgHandle, frameIdx); exists {
		anchor := GetAnchor(entry)

		width, height := img.Bounds().Dx(), img.Bounds().Dy()
		bounds.Min = geom.Vec2{
			X: -anchor.X * float64(width),
			Y: -anchor.Y * float64(height),
		}
		bounds.Max = geom.Vec2{
			X: bounds.Min.X + float64(width),
			Y: bounds.Min.Y + float64(height),
		}
	}

	transform.AddTransform(entry,
		transform.WithBounds(bounds.Min, bounds.Max),
	)

	renderer.AddRenderer(entry,
		renderer.WithRendererType(p.rendererType),
	)

	return entry
}
