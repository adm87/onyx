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
	renderingTasks []engine.RenderingTask
	rendererTypes  []rendering.RendererType
}

func NewImageRenderingAdapter(assetAdapter *ImageAssetAdapter) *ImageRenderingAdapter {
	return &ImageRenderingAdapter{
		assetAdapter:   assetAdapter,
		renderingTasks: make([]engine.RenderingTask, 0, 100),
		rendererTypes:  []rendering.RendererType{ImageRendererType},
	}
}

func (a *ImageRenderingAdapter) SupportedRendererTypes() []rendering.RendererType {
	return a.rendererTypes
}

func (a *ImageRenderingAdapter) GetRenderingTasks(entry *donburi.Entry, viewport geom.AABB, viewMatrix ebiten.GeoM) []engine.RenderingTask {
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
	scale := transform.GetScale(entry)

	aX := anchor.X * float64(img.Bounds().Dx()) * scale.X
	aY := anchor.Y * float64(img.Bounds().Dy()) * scale.Y

	matrix.Translate(-aX, -aY)
	matrix.Concat(viewMatrix)

	opts := ebiten.DrawImageOptions{
		Filter: filter,
	}

	alpha := float32(color.A) / 255

	a.renderingTasks = append(a.renderingTasks, engine.RenderingTask{
		Render: func(screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
			opts.GeoM = matrix

			// NOTE: Just using ScaleWithColor doesn't work as expected, it seems to ignore the alpha component of the color.
			opts.ColorScale.ScaleWithColor(color)
			opts.ColorScale.ScaleAlpha(alpha)

			screen.DrawImage(img, &opts)
			return nil
		},
		Layer:  rendering.GetLayer(entry),
		ZIndex: rendering.GetZIndex(entry),
	})

	return a.renderingTasks
}
