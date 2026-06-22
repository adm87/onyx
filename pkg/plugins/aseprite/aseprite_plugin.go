package aseprite

import (
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/yohamta/donburi"
)

type AsepritePlugin struct {
	imagePlugin *images.ImagePlugin
	library     *AsepriteLibrary
	systems     *AsepriteSystems
}

func NewAsepritePlugin(imagePlugin *images.ImagePlugin) *AsepritePlugin {
	library := NewAsepriteLibrary(imagePlugin.Assets())
	return &AsepritePlugin{
		imagePlugin: imagePlugin,
		library:     library,
		systems:     NewAsepriteSystems(library),
	}
}

func (a *AsepritePlugin) Library() *AsepriteLibrary {
	return a.library
}

func (a *AsepritePlugin) Systems() *AsepriteSystems {
	return a.systems
}

func (a *AsepritePlugin) CreateSprite(ecs donburi.World, opts ...SpriteOption) *donburi.Entry {
	options := DefaultSpriteOptions()
	for _, opt := range opts {
		opt(options)
	}

	entry := a.imagePlugin.CreateImage(ecs,
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
