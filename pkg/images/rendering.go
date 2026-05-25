package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func NewRenderingSystem() func(world donburi.World) []engine.RenderTask {
	var tasks []engine.RenderTask
	return func(world donburi.World) []engine.RenderTask {
		tasks = tasks[:0]
		query.Each(world, func(e *donburi.Entry) {
			img := rendering.GetImage(e)
			if img == nil {
				return // Don't enqueue render tasks for entities without an image
			}
			vis := rendering.IsVisible(e)
			if !vis {
				return // Don't enqueue render tasks for invisible entities
			}

			layer := rendering.GetLayer(e)
			zIndex := rendering.GetZIndex(e)
			matrix := transform.GetMatrix(e)
			anchor := rendering.GetAnchor(e)
			col := rendering.GetColor(e)

			aX := anchor.X * float64(img.Bounds().Dx())
			aY := anchor.Y * float64(img.Bounds().Dy())

			tasks = append(tasks, engine.RenderTask{
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
		return tasks
	}
}
