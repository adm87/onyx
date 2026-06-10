package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type renderingAdapter struct {
	assetAdapter   *assetAdapter
	renderingTasks []engine.RenderingTask
}

func newRendererAdapter(assetAdapter *assetAdapter) *renderingAdapter {
	return &renderingAdapter{
		assetAdapter:   assetAdapter,
		renderingTasks: make([]engine.RenderingTask, 1),
	}
}

func (a *renderingAdapter) GetRenderingTasks(entry *donburi.Entry, viewport geom.AABB, viewMatrix ebiten.GeoM) []engine.RenderingTask {
	a.renderingTasks = a.renderingTasks[:0]

	visible := rendering.IsVisible(entry)
	if !visible {
		return a.renderingTasks
	}

	layer := rendering.GetLayer(entry)
	zIndex := rendering.GetZIndex(entry)
	color := rendering.GetColor(entry)
	filter := rendering.GetFilter(entry)
	anchor := rendering.GetAnchor(entry)

	img, exists := a.assetAdapter.store.Get(GetImageHandle(entry))
	assert.True(exists, "cannot find image")

	scale := transform.GetScale(entry)

	// TODO - revist this, should it be scale or sign of scale?
	aX := float64(img.Bounds().Dx()) * anchor.X * scale.X
	aY := float64(img.Bounds().Dy()) * anchor.Y * scale.Y

	matrix := transform.GetMatrix(entry)
	matrix.Translate(-aX, -aY)
	matrix.Concat(viewMatrix)

	a.renderingTasks = append(a.renderingTasks, engine.RenderingTask{
		Layer:  layer,
		ZIndex: zIndex,
		Job: func(target *ebiten.Image) {
			opt := ebiten.DrawImageOptions{
				Filter: filter,
				GeoM:   matrix,
			}

			opt.ColorScale.ScaleWithColor(color)
			opt.ColorScale.ScaleAlpha(float32(color.A) / 255)

			target.DrawImage(img, &opt)
		},
	})

	return a.renderingTasks
}
