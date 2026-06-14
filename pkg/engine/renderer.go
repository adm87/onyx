package engine

import (
	"slices"

	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
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
		entry *donburi.Entry,
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

	jobs     []*RenderingJob
	entries  []*donburi.Entry
	entities []donburi.Entity
}

func newRenderer(logger Logger) *renderer {
	return &renderer{
		logger:   logger,
		adapters: slotmap.New[RenderingAdapter](0),
		jobPool: &renderingJobPool{
			pool: make([]*RenderingJob, 0, 100),
		},
		jobs:     make([]*RenderingJob, 0, 100),
		entries:  make([]*donburi.Entry, 0, 100),
		entities: make([]donburi.Entity, 0, 100),
	}
}

func (r *renderer) AddRenderingAdapter(adapter RenderingAdapter) uint64 {
	return r.adapters.Insert(adapter)
}

func (r *renderer) render(ecs donburi.World, world *world, screen *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) {
	r.entities = world.entities.QueryInto(viewport, r.entities[:0])

	r.entries = r.entries[:0]
	r.jobs = r.jobs[:0]
	r.jobPool.i = 0

	for i := range r.entities {
		entry := ecs.Entry(r.entities[i])

		renderer := rendering.GetRenderer(entry)
		if !renderer.Visible {
			continue
		}

		aabb := transform.GetWorldBounds(entry)
		if !aabb.Intersects(viewport) {
			continue
		}

		rendererID := rendering.GetRendererID(entry)
		adapter, exists := r.adapters.Get(rendererID)

		if !exists {
			r.logger.Warn("renderer adapter not found for entry renderer: %v", rendererID)
			continue
		}

		jobs := adapter.GetJobs(entry, viewport, viewMatrix, r.jobPool)

		r.jobs = append(r.jobs, jobs...)
		r.entries = append(r.entries, entry)
	}

	slices.SortFunc(r.jobs, func(a, b *RenderingJob) int {
		if a.Layer == b.Layer {
			return a.ZIndex - b.ZIndex
		}
		return a.Layer - b.Layer
	})

	for _, job := range r.jobs {
		if job.Buffer != nil {
			screen.DrawImage(job.Buffer, &job.Options)
		}
	}
}
