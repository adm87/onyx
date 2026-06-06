package engine

import (
	"cmp"
	"slices"

	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/spatialhash"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RenderingAdapter interface {
	GetRenderTasks(entry *donburi.Entry, layer int, zIndex int, viewMatrix ebiten.GeoM) []RenderTask
	SupportedRendererTypes() []rendering.RendererType
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

type renderableIndexing map[donburi.Entity]spatialhash.SpatialIndex

type renderer struct {
	logger *logger
	queue  []RenderTask

	adapters  map[AdapterID]RenderingAdapter
	renderers map[rendering.RendererType]RenderingAdapter

	partitioning *spatialhash.SpatialHash[donburi.Entity]
	indexing     renderableIndexing
}

func newRenderer(logger *logger) *renderer {
	return &renderer{
		logger:    logger,
		queue:     make([]RenderTask, 0, 100),
		adapters:  make(map[AdapterID]RenderingAdapter),
		renderers: make(map[rendering.RendererType]RenderingAdapter),
		partitioning: spatialhash.New[donburi.Entity](
			spatialhash.WithResolutions(16, 32),
			spatialhash.WithCapacity(100),
		),
		indexing: make(renderableIndexing),
	}
}

func (r *renderer) addRenderable(entry *donburi.Entry, aabb geom.AABB) bool {
	entity := entry.Entity()

	if _, exists := r.indexing[entity]; exists {
		return false // Entity is already indexed, cannot add again
	}

	index, ok := r.partitioning.Insert(entity, aabb)
	if ok {
		r.indexing[entity] = index
	}

	return ok
}

func (r *renderer) removeRenderable(entry *donburi.Entry) bool {
	entity := entry.Entity()

	index, exists := r.indexing[entity]
	if !exists {
		return false // Entity is not indexed, cannot remove
	}

	ok := r.partitioning.Remove(index)
	if ok {
		delete(r.indexing, entity)
	}

	return ok
}

func (r *renderer) updateRenderable(entry *donburi.Entry, aabb geom.AABB) bool {
	entity := entry.Entity()

	index, exists := r.indexing[entity]
	if !exists {
		return false // Entity is not indexed, cannot update
	}

	return r.partitioning.Reinsert(index, aabb)
}

func (r *renderer) render(ecs donburi.World, screen *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error {
	r.queue = r.queue[:0]

	r.partitioning.QueryAll(viewport, func(e donburi.Entity) bool {
		entry := ecs.Entry(e)

		visible := rendering.IsVisible(entry)
		if !visible {
			return true // Skip invisible entities
		}

		layer := rendering.GetLayer(entry)
		zIndex := rendering.GetZIndex(entry)
		rendererType := rendering.GetRendererType(entry)

		adapter, exists := r.renderers[rendererType]
		if !exists {
			return true // No adapter for this renderer type, skip
		}

		tasks := adapter.GetRenderTasks(entry, layer, zIndex, viewMatrix)
		r.queue = append(r.queue, tasks...)

		return true
	})

	println("Render queue length:", len(r.queue))

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

	for _, rendererType := range adapter.SupportedRendererTypes() {
		if _, exists := r.renderers[rendererType]; exists {
			r.logger.Warn("Renderer for type '%s' already exists, skipping", rendererType)
			continue
		}
		r.renderers[rendererType] = adapter
	}

	r.adapters[adapterID] = adapter
}

func (r *renderer) GetRenderingAdapter(adapterID AdapterID) (RenderingAdapter, bool) {
	adapter, found := r.adapters[adapterID]
	return adapter, found
}
