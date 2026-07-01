package aseprite

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/yohamta/donburi"
)

// pluginID is a unique identifier for the AsepritePlugin type.
var pluginID = engine.TypeHash[AsepritePlugin]()

func PluginID() uint64 {
	return pluginID
}

type AsepritePlugin interface {
	engine.Plugin

	Library() *AsepriteLibrary
	Systems() *AsepriteSystems

	CreateSprite(ecs donburi.World, opts ...SpriteOption) *donburi.Entry
}

type plugin struct {
	library *AsepriteLibrary
	systems *AsepriteSystems

	imagePlugin images.ImagePlugin
}

func NewPlugin() AsepritePlugin {
	library := NewAsepriteLibrary()
	systems := NewAsepriteSystems(library)
	return &plugin{
		library: library,
		systems: systems,
	}
}

func (p *plugin) OnRegister(game engine.Game) {
	imagePlugin := engine.GetPlugin[images.ImagePlugin](game, images.PluginID())
	p.imagePlugin = imagePlugin
	p.library.imageAssets = imagePlugin.Assets()
}

func (p *plugin) ID() uint64 {
	return PluginID()
}

func (p *plugin) Library() *AsepriteLibrary {
	return p.library
}

func (p *plugin) Systems() *AsepriteSystems {
	return p.systems
}

func (p *plugin) CreateSprite(ecs donburi.World, opts ...SpriteOption) *donburi.Entry {
	options := DefaultSpriteOptions()
	for _, opt := range opts {
		opt(options)
	}

	entry := p.imagePlugin.CreateImage(ecs,
		images.WithHandle(options.ImageOptions.Handle),
		images.WithAnchor(options.ImageOptions.Anchor.X, options.ImageOptions.Anchor.Y),
		images.WithFrame(options.ImageOptions.Frame),
		images.WithColor(options.ImageOptions.Color),
		images.WithFilter(options.ImageOptions.Filter),
	)

	SetAnimationState(entry, options.State)
	SetAnimator(entry, &AnimatorModel{
		Loops:     options.Loops,
		direction: 1,
	})
	SetClip(entry, options.Clip)

	return entry
}
