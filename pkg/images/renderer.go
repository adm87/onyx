package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/asset"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ImageRenderingAdapter struct {
	assetAdapter   *ImageAssetAdapter
	renderingTasks []engine.RenderTask
	rendererTypes  []rendering.RendererType
}

func NewImageRenderingAdapter(assetAdapter *ImageAssetAdapter) *ImageRenderingAdapter {
	return &ImageRenderingAdapter{
		assetAdapter:   assetAdapter,
		renderingTasks: make([]engine.RenderTask, 0, 100),
		rendererTypes:  []rendering.RendererType{ImageRendererType},
	}
}

func (a *ImageRenderingAdapter) SupportedRendererTypes() []rendering.RendererType {
	return a.rendererTypes
}

func (a *ImageRenderingAdapter) GetRenderTasks(entry *donburi.Entry, layer int, zIndex int, viewport geom.AABB, viewMatrix ebiten.GeoM) []engine.RenderTask {
	a.renderingTasks = a.renderingTasks[:0]

	ref := asset.GetAssetReference(entry)
	if ref == asset.UnknownRef {
		return a.renderingTasks // Don't enqueue render tasks for entities without an image reference
	}

	img, exists := a.assetAdapter.GetImage(ref)
	if !exists {
		return a.renderingTasks // Don't enqueue render tasks for entities with an invalid image reference
	}

	matrix := transform.GetMatrix(entry)
	anchor := rendering.GetAnchor(entry)
	color := rendering.GetColor(entry)
	filter := rendering.GetFilter(entry)

	aX := anchor.X * float64(img.Bounds().Dx())
	aY := anchor.Y * float64(img.Bounds().Dy())

	opts := ebiten.DrawImageOptions{
		Filter: filter,
	}

	alpha := float32(color.A) / 255

	a.renderingTasks = append(a.renderingTasks, engine.RenderTask{
		Render: func(screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
			opts.GeoM.Reset()
			opts.GeoM.Translate(-aX, -aY)
			opts.GeoM.Concat(matrix)
			opts.GeoM.Concat(viewMatrix)

			// NOTE: Just using ScaleWithColor doesn't work as expected, it seems to ignore the alpha component of the color.
			opts.ColorScale.ScaleWithColor(color)
			opts.ColorScale.ScaleAlpha(alpha)

			screen.DrawImage(img, &opts)
			return nil
		},
		Layer:  layer,
		ZIndex: zIndex,
	})

	return a.renderingTasks
}
