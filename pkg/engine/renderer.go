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

type renderTask struct {
	render func(screen *ebiten.Image) error
	layer  int
	zIndex int
}

type renderer struct {
	logger     *logger
	imageQuery *donburi.Query

	queue []renderTask
}

func newRenderer(logger *logger) *renderer {
	return &renderer{
		logger: logger,
		queue:  make([]renderTask, 0, 100),
		imageQuery: donburi.NewQuery(
			filter.Contains(transform.Matrix, rendering.Renderer, rendering.Image),
		),
	}
}

func (r *renderer) render(world donburi.World, screen *ebiten.Image) error {
	r.buildRenderQueue(world)

	slices.SortFunc(r.queue, func(a, b renderTask) int {
		if a.layer != b.layer {
			return cmp.Compare(a.layer, b.layer)
		}
		return cmp.Compare(a.zIndex, b.zIndex)
	})

	for _, task := range r.queue {
		if err := task.render(screen); err != nil {
			return err
		}
	}

	return nil
}

func (r *renderer) buildRenderQueue(world donburi.World) {
	r.queue = r.queue[:0]

	r.imageQuery.Each(world, func(e *donburi.Entry) {
		img := rendering.GetImage(e)
		if img == nil {
			return // Don't enqueue render tasks for entities without an image
		}
		vis := rendering.IsVisible(e)
		if !vis {
			return // Don't enqueue render tasks for invisible entities
		}

		layer := rendering.GetLayer(e)
		zIndex := rendering.GetZIndex(e)
		matrix := transform.GetMatrix(e)
		anchor := rendering.GetAnchor(e)
		color := rendering.GetColor(e)

		aX := anchor.X * float64(img.Bounds().Dx())
		aY := anchor.Y * float64(img.Bounds().Dy())

		r.queue = append(r.queue, renderTask{
			render: func(screen *ebiten.Image) error {
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(-aX, -aY)
				opts.GeoM.Concat(matrix)
				opts.ColorScale.ScaleWithColor(color)
				screen.DrawImage(img, opts)
				return nil
			},
			layer:  layer,
			zIndex: zIndex,
		})
	})
}
