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
	renderer *rendering.RendererModel,
	bounds geom.AABB,
	viewport geom.AABB,
	viewMatrix ebiten.GeoM,
	pool engine.RenderingJobPool) []*engine.RenderingJob {
	a.jobs = a.jobs[:0]

	image := GetImage(entry)
	trans := transform.GetTransform(entry)
	matrix := transform.GetMatrix(entry)

	img, exists := a.assetAdapter.getFrame(image.Handle, image.Frame)
	if !exists {
		return a.jobs
	}

	aX := float64(img.Bounds().Dx()) * image.Anchor.X * trans.ScaleX
	aY := float64(img.Bounds().Dy()) * image.Anchor.Y * trans.ScaleY

	matrix.Translate(-aX, -aY)
	matrix.Concat(viewMatrix)

	job := pool.Get()
	job.Buffer = img
	job.Layer = renderer.Layer
	job.ZIndex = renderer.ZIndex
	job.Options.Filter = image.Filter
	job.Options.GeoM = matrix
	job.Options.ColorScale.ScaleWithColor(image.Color)
	job.Options.ColorScale.ScaleAlpha(float32(image.Color.A) / 255)
	a.jobs = append(a.jobs, job)

	return a.jobs
}
