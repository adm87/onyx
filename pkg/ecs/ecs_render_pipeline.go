package ecs

import (
	"slices"

	"github.com/adm87/onyx/pkg/ecs/camera"
	"github.com/adm87/onyx/pkg/ecs/renderer"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ECSRenderAdapter interface {
	PrepareRenderingTasks(
		entry *donburi.Entry,
		renderer *renderer.RendererModel,
		pool *engine.RenderingPool,
		viewport geom.AABB,
		viewMatrix ebiten.GeoM) []*engine.RenderingTask
	SetAdapterIndex(id uint64)
}

type ECSRenderPipeline struct {
	world donburi.World

	screen engine.Screen
	logger engine.Logger

	adapters    *slotmap.SlotMap[ECSRenderAdapter]
	partitioner *ECSPartitioner

	tasks []*engine.RenderingTask
}

func NewECSRenderPipeline(
	world donburi.World,
	screen engine.Screen,
	logger engine.Logger,
	partitioner *ECSPartitioner) *ECSRenderPipeline {
	return &ECSRenderPipeline{
		world:       world,
		screen:      screen,
		logger:      logger,
		adapters:    slotmap.New[ECSRenderAdapter](0),
		partitioner: partitioner,
		tasks:       make([]*engine.RenderingTask, 0, 100),
	}
}

func (r *ECSRenderPipeline) AddAdapters(adapters ...ECSRenderAdapter) *ECSRenderPipeline {
	for _, adapter := range adapters {
		index := r.adapters.Insert(adapter)
		adapter.SetAdapterIndex(index)
	}
	return r
}

func (r *ECSRenderPipeline) GetRenderingTasks(pool *engine.RenderingPool) []*engine.RenderingTask {
	mainCamera, found := camera.GetMainCamera(r.world)
	if !found {
		r.logger.Warn("No main camera found in the world.")
		return nil
	}

	safeArea := r.screen.SafeArea()

	viewport := camera.GetViewport(mainCamera, safeArea)
	viewMatrix := camera.GetViewMatrix(mainCamera, safeArea)

	r.tasks = r.tasks[:0]
	r.partitioner.entities.Query(viewport, func(item donburi.Entity) {
		entry := r.world.Entry(item)

		renderer := renderer.GetRenderer(entry)
		if renderer == nil || !renderer.Visible {
			return
		}

		if adapter, exists := r.adapters.Get(renderer.Type); exists {
			tasks := adapter.PrepareRenderingTasks(
				entry,
				renderer,
				pool,
				viewport,
				viewMatrix,
			)
			r.tasks = append(r.tasks, tasks...)
		}
	})

	slices.SortFunc(r.tasks, func(i, j *engine.RenderingTask) int {
		if i.Layer != j.Layer {
			return i.Layer - j.Layer
		}
		return i.ZIndex - j.ZIndex
	})

	return r.tasks
}
