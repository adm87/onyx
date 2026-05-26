package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/images/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ImageRenderingAdapter struct {
	renderingTasks []engine.RenderTask
}

func NewImageRenderingAdapter() *ImageRenderingAdapter {
	return &ImageRenderingAdapter{
		renderingTasks: make([]engine.RenderTask, 0, 100),
	}
}

func (a *ImageRenderingAdapter) GetRenderTasks(world donburi.World) []engine.RenderTask {
	a.renderingTasks = a.renderingTasks[:0]

	ImageQuery.Each(world, func(e *donburi.Entry) {
		img := components.GetImage(e)
		if img == nil {
			return // Don't enqueue render tasks for entities without an image
		}

		vis := rendering.IsVisible(e)
		if !vis {
			return // Don't enqueue render tasks for invisible entities
		}

		layer := rendering.GetLayer(e)
		zIndex := rendering.GetZIndex(e)
		anchor := rendering.GetAnchor(e)
		col := rendering.GetColor(e)

		matrix := transform.GetMatrix(e)

		aX := anchor.X * float64(img.Bounds().Dx())
		aY := anchor.Y * float64(img.Bounds().Dy())

		a.renderingTasks = append(a.renderingTasks, engine.RenderTask{
			Render: func(screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
				opts := ebiten.DrawImageOptions{}
				opts.GeoM.Translate(-aX, -aY)
				opts.GeoM.Concat(matrix)
				opts.GeoM.Concat(viewMatrix)
				opts.ColorScale.ScaleWithColor(col)
				opts.ColorScale.ScaleAlpha(float32(col.A) / 255)
				screen.DrawImage(img, &opts)
				return nil
			},
			Layer:  layer,
			ZIndex: zIndex,
		})
	})

	return a.renderingTasks
}
