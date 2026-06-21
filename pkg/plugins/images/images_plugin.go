package images

import (
	"github.com/adm87/onyx/pkg/ecs/renderer"
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
)

type ImagePlugin struct {
	assets   *ImageAssets
	renderer *ImageECSRenderer
}

func NewImagePlugin() *ImagePlugin {
	assets := NewImageAssets()
	return &ImagePlugin{
		assets:   assets,
		renderer: NewImageECSRenderer(assets),
	}
}

func (i *ImagePlugin) Assets() *ImageAssets {
	return i.assets
}

func (i *ImagePlugin) Renderer() *ImageECSRenderer {
	return i.renderer
}

func (i *ImagePlugin) CreateImage(world donburi.World, opts ...Option) *donburi.Entry {
	entry := NewImage(world, opts...)

	var bounds geom.AABB

	imgHandle := GetHandle(entry)
	frameIdx := GetFrame(entry)

	if img, exists := i.assets.GetFrame(imgHandle, frameIdx); exists {
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
		renderer.WithRendererType(i.renderer.adapterID),
	)

	return entry
}
