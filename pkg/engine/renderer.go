package engine

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
	"github.com/hajimehoshi/ebiten/v2"
)

type RenderingJob struct {
	Layer   int
	ZIndex  int
	Buffer  *ebiten.Image
	Options ebiten.DrawImageOptions
}

type RenderingJobPool interface {
	Get() *RenderingJob
}

type RenderingAdapter interface {
	GetJobs(
		bounds geom.AABB,
		viewport geom.AABB,
		viewMatrix ebiten.GeoM,
		pool RenderingJobPool) []*RenderingJob
}

type Renderer interface {
	AddRenderingAdapter(adapter RenderingAdapter) uint64
}

type renderingJobPool struct {
	pool []*RenderingJob
	i    int
}

func (p *renderingJobPool) Get() *RenderingJob {
	if p.i >= len(p.pool) {
		p.pool = append(p.pool, &RenderingJob{})
	}
	job := p.pool[p.i]
	job.Buffer = nil
	job.Options.GeoM.Reset()
	job.Options.ColorScale.Reset()
	p.i++
	return job
}

type renderer struct {
	logger Logger

	adapters *slotmap.SlotMap[RenderingAdapter]
	jobPool  *renderingJobPool

	jobs []*RenderingJob
}

func newRenderer(logger Logger) *renderer {
	return &renderer{
		logger:   logger,
		adapters: slotmap.New[RenderingAdapter](0),
		jobPool: &renderingJobPool{
			pool: make([]*RenderingJob, 0, 100),
		},
		jobs: make([]*RenderingJob, 0, 100),
	}
}

func (r *renderer) AddRenderingAdapter(adapter RenderingAdapter) uint64 {
	return r.adapters.Insert(adapter)
}

func (r *renderer) render(screen *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) {
	r.jobs = r.jobs[:0]
	r.jobPool.i = 0
}
