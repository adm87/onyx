package engine

import (
	"cmp"
	"slices"

	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type RenderTask struct {
	Render func(screen *ebiten.Image, viewMatrix ebiten.GeoM) error
	Layer  int
	ZIndex int
}

type Renderer interface {
	AddRenderingSystem(system func(world donburi.World) []RenderTask)
}

type renderer struct {
	logger     *logger
	imageQuery *donburi.Query

	systems []func(world donburi.World) []RenderTask
	queue   []RenderTask
}

func newRenderer(logger *logger) *renderer {
	return &renderer{
		logger: logger,
		queue:  make([]RenderTask, 0, 100),
		imageQuery: donburi.NewQuery(
			filter.Contains(transform.Matrix, rendering.Renderer, rendering.Image),
		),
	}
}

func (r *renderer) render(world donburi.World, screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
	r.queue = r.queue[:0]

	for _, system := range r.systems {
		r.queue = append(r.queue, system(world)...)
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

func (r *renderer) AddRenderingSystem(system func(world donburi.World) []RenderTask) {
	r.systems = append(r.systems, system)
}
