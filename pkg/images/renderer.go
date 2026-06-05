package images

import (
	"image/color"

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
}

func NewImageRenderingAdapter(assetAdapter *ImageAssetAdapter) *ImageRenderingAdapter {
	return &ImageRenderingAdapter{
		assetAdapter:   assetAdapter,
		renderingTasks: make([]engine.RenderTask, 0, 100),
	}
}

func (a *ImageRenderingAdapter) GetRenderTasks(world donburi.World, viewMatrix ebiten.GeoM) []engine.RenderTask {
	a.renderingTasks = a.renderingTasks[:0]

	rendering.QueryVisibleWith(world, ImageQuery,
		func(entry *donburi.Entry, anchor geom.Vec2, color color.RGBA, filter ebiten.Filter, visible bool, layer, zIndex int) {
			ref := asset.GetAssetReference(entry)
			if ref == asset.UnknownRef {
				return // Don't enqueue render tasks for entities without an image reference
			}

			img, exists := a.assetAdapter.GetImage(ref)
			if !exists {
				return // Don't enqueue render tasks for entities with an invalid image reference
			}

			matrix := transform.GetMatrix(entry)

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
		},
	)

	return a.renderingTasks
}
