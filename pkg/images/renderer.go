package images

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type renderingAdapter struct {
	assetAdapter *assetAdapter
	jobs         []*engine.RenderingJob
}

func newRendererAdapter(assetAdapter *assetAdapter) *renderingAdapter {
	return &renderingAdapter{
		assetAdapter: assetAdapter,
		jobs:         make([]*engine.RenderingJob, 0, 100),
	}
}

func (a *renderingAdapter) GetJobs(
	entry *donburi.Entry,
	viewport geom.AABB,
	viewMatrix ebiten.GeoM,
	pool engine.RenderingJobPool) []*engine.RenderingJob {
	a.jobs = a.jobs[:0]

	handle, exists := GetImageHandle(entry)
	if !exists {
		return nil
	}

	layer := rendering.GetLayer(entry)
	zIndex := rendering.GetZIndex(entry)
	color := rendering.GetColor(entry)
	filter := rendering.GetFilter(entry)
	anchor := rendering.GetAnchor(entry)

	img, exists := a.assetAdapter.store.Get(handle)
	assert.True(exists, "cannot find image")

	scale := transform.GetScale(entry)

	// TODO - revist this, should it be scale or sign of scale?
	aX := float64(img.Bounds().Dx()) * anchor.X * scale.X
	aY := float64(img.Bounds().Dy()) * anchor.Y * scale.Y

	matrix := transform.GetMatrix(entry)
	matrix.Translate(-aX, -aY)
	matrix.Concat(viewMatrix)

	job := pool.Get(img)
	job.Layer = layer
	job.ZIndex = zIndex
	job.Options.Filter = filter
	job.Options.GeoM = matrix
	job.Options.ColorScale.ScaleWithColor(color)
	job.Options.ColorScale.ScaleAlpha(float32(color.A) / 255)

	a.jobs = append(a.jobs, job)
	return a.jobs
}

func (a *renderingAdapter) drawImage(
	img *ebiten.Image,
	matrix ebiten.GeoM,
	color color.RGBA,
	filter ebiten.Filter) func(target *ebiten.Image) {
	return func(target *ebiten.Image) {
		opt := ebiten.DrawImageOptions{
			Filter: filter,
			GeoM:   matrix,
		}

		opt.ColorScale.ScaleWithColor(color)
		opt.ColorScale.ScaleAlpha(float32(color.A) / 255)

		target.DrawImage(img, &opt)
	}
}
