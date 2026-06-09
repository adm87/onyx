package engine

import (
	"slices"

	"github.com/adm87/onyx/pkg/assert"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/spatialhash"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RenderingJob func(target *ebiten.Image)

type RenderingTask struct {
	Layer  int
	ZIndex int
	Job    RenderingJob
}

type RenderingAdapter interface {
	GetRenderingTasks(entry *donburi.Entry, viewport geom.AABB, viewMatrix ebiten.GeoM) []RenderingTask
}

type Renderer interface {
	AddRenderingAdapter(adapter RenderingAdapter) uint64
}

type renderer struct {
	adapters  *slotmap.SlotMap[RenderingAdapter]
	partition *spatialhash.SpatialHash[donburi.Entity]

	tasks []RenderingTask
}

func newRenderer() *renderer {
	return &renderer{
		adapters: slotmap.New[RenderingAdapter](0),
		partition: spatialhash.New[donburi.Entity](64,
			spatialhash.Padding{Left: 1, Right: 1},
		),
		tasks: make([]RenderingTask, 0, 100),
	}
}

func (r *renderer) AddRenderingAdapter(adapter RenderingAdapter) uint64 {
	return r.adapters.Insert(adapter)
}

func (r *renderer) render(ecs donburi.World, screen *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) {
	r.tasks = r.tasks[:0]

	queryRegion := viewport.Scale(2)
	r.partition.Query(queryRegion, func(entity donburi.Entity) {
		entry := ecs.Entry(entity)

		adapter, exists := r.adapters.Get(rendering.GetRenderer(entry))
		assert.True(exists, "cannot find rendering adapter")

		tasks := adapter.GetRenderingTasks(entry, viewport, viewMatrix)
		r.tasks = append(r.tasks, tasks...)
	})

	slices.SortFunc(r.tasks, func(a, b RenderingTask) int {
		if a.Layer == b.Layer {
			return a.ZIndex - b.ZIndex
		}
		return a.Layer - b.Layer
	})

	for _, task := range r.tasks {
		task.Job(screen)
	}
}
