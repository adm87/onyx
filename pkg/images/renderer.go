package images

import (
	"github.com/adm87/onyx/pkg/engine"
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

func newRenderingAdapter(assetAdapter *assetAdapter) *renderingAdapter {
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
		return a.jobs
	}

	layer := rendering.GetLayer(entry)
	zIndex := rendering.GetZIndex(entry)
	color := rendering.GetColor(entry)
	filter := rendering.GetFilter(entry)
	anchor := rendering.GetAnchor(entry)

	img, exists := a.assetAdapter.store.Get(handle)
	if !exists {
		return a.jobs
	}

	// TODO - revist this, should it be scale or sign of scale?
	scale := transform.GetScale(entry)
	aX := float64(img.Bounds().Dx()) * anchor.X * scale.X
	aY := float64(img.Bounds().Dy()) * anchor.Y * scale.Y

	matrix := transform.GetMatrix(entry)
	matrix.Translate(-aX, -aY)
	matrix.Concat(viewMatrix)

	job := pool.Get()
	job.Buffer = img
	job.Layer = layer
	job.ZIndex = zIndex
	job.Options.Filter = filter
	job.Options.GeoM = matrix
	job.Options.ColorScale.ScaleWithColor(color)
	job.Options.ColorScale.ScaleAlpha(float32(color.A) / 255)
	a.jobs = append(a.jobs, job)

	return a.jobs
}
