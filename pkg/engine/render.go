package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer interface {
	UsePipeline(RenderPipeline)
}

type RenderPipeline interface {
	GetRenderingTasks(taskPool *RenderingPool) []*RenderingTask
}

type RenderingTask struct {
	Buffer  *ebiten.Image
	Options *ebiten.DrawImageOptions
	Layer   int
	ZIndex  int
}

type RenderingPool struct {
	pool []*RenderingTask
	i    int
}

func (p *RenderingPool) Get() *RenderingTask {
	if p.i >= len(p.pool) {
		p.pool = append(p.pool, &RenderingTask{})
	}
	task := p.pool[p.i]
	task.Buffer = nil
	task.Options = nil
	p.i++
	return task
}

type renderer struct {
	logger *logger

	pool  *RenderingPool
	tasks []*RenderingTask

	pipeline RenderPipeline
}

func newRenderer(logger *logger) *renderer {
	return &renderer{
		logger: logger,
		pool: &RenderingPool{
			pool: make([]*RenderingTask, 0, 100),
		},
		tasks: make([]*RenderingTask, 0, 100),
	}
}

func (r *renderer) UsePipeline(p RenderPipeline) {
	r.pipeline = p
}

func (r *renderer) render(target *ebiten.Image) {
	if r.pipeline == nil {
		r.logger.Warn("No render pipeline set. Skipping rendering.")
		return
	}

	r.tasks = r.tasks[:0]
	r.pool.i = 0

	r.tasks = append(r.tasks, r.pipeline.GetRenderingTasks(r.pool)...)
	for i := range r.tasks {
		target.DrawImage(r.tasks[i].Buffer, r.tasks[i].Options)
	}
}
