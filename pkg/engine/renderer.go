package engine

import (
	"cmp"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RenderingAdapter interface {
	GetRenderTasks(ecs donburi.World, viewMatrix ebiten.GeoM) []RenderTask
}

type RenderTask struct {
	Render func(screen *ebiten.Image, viewMatrix ebiten.GeoM) error
	Layer  int
	ZIndex int
}

type Renderer interface {
	AddRenderingAdapter(adapterID AdapterID, adapter RenderingAdapter)
	GetRenderingAdapter(adapterID AdapterID) (RenderingAdapter, bool)
}

type renderer struct {
	logger   *logger
	queue    []RenderTask
	adapters map[AdapterID]RenderingAdapter
}

func newRenderer(logger *logger) *renderer {
	return &renderer{
		logger:   logger,
		queue:    make([]RenderTask, 0, 100),
		adapters: make(map[AdapterID]RenderingAdapter),
	}
}

func (r *renderer) render(ecs donburi.World, screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
	r.queue = r.queue[:0]

	for _, adapter := range r.adapters {
		r.queue = append(r.queue, adapter.GetRenderTasks(ecs, viewMatrix)...)
	}

	slices.SortFunc(r.queue, func(a, b RenderTask) int {
		if a.Layer != b.Layer {
			return cmp.Compare(a.Layer, b.Layer)
		}
		return cmp.Compare(a.ZIndex, b.ZIndex)
	})

	for _, task := range r.queue {
		if err := task.Render(screen, viewMatrix); err != nil {
			return err
		}
	}

	return nil
}

func (r *renderer) AddRenderingAdapter(adapterID AdapterID, adapter RenderingAdapter) {
	if _, exists := r.adapters[adapterID]; exists {
		r.logger.Warn("Rendering adapter with ID '%s' already exists, skipping", adapterID)
		return
	}

	r.adapters[adapterID] = adapter
}

func (r *renderer) GetRenderingAdapter(adapterID AdapterID) (RenderingAdapter, bool) {
	adapter, found := r.adapters[adapterID]
	return adapter, found
}
