package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/spatialhash"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type World interface {
	ECS() donburi.World

	Add(entry *donburi.Entry)
	Remove(entry *donburi.Entry)
	Update(entry *donburi.Entry)

	QueryRegion(region geom.AABB, callback func(*donburi.Entry))
}

type world struct {
	ecs donburi.World

	renderer *renderer

	entities     *spatialhash.SpatialHash[donburi.Entity]
	queryResults []*donburi.Entry
}

var worldIndexing = donburi.NewComponentType[uint64]()

func newWorld(renderer *renderer) *world {
	ecs := donburi.NewWorld()
	return &world{
		ecs:          ecs,
		renderer:     renderer,
		entities:     spatialhash.New[donburi.Entity](16, spatialhash.Padding{}),
		queryResults: make([]*donburi.Entry, 0, 100),
	}
}

func (w *world) ECS() donburi.World {
	return w.ecs
}

func (w *world) Add(entry *donburi.Entry) {
	aabb := transform.GetLocalBounds(entry).Translate(transform.GetPosition(entry))
	entity := entry.Entity()
	index := w.entities.Insert(entity, aabb)
	donburi.Add(entry, worldIndexing, &index)
}

func (w *world) Remove(entry *donburi.Entry) {
	index := worldIndexing.Get(entry)
	w.entities.Remove(*index)
	entry.Remove()
}

func (w *world) Update(entry *donburi.Entry) {
	index := worldIndexing.Get(entry)
	aabb := transform.GetLocalBounds(entry).Translate(transform.GetPosition(entry))
	w.entities.Update(*index, aabb)
}

func (w *world) UpdateBounds(entry *donburi.Entry, bounds geom.AABB) {
	index := worldIndexing.Get(entry)
	w.entities.Update(*index, bounds)
}

func (w *world) QueryRegion(region geom.AABB, callback func(*donburi.Entry)) {
	w.entities.Query(region, func(entity donburi.Entity) {
		entry := w.ecs.Entry(entity)
		aabb := transform.GetLocalBounds(entry).Translate(transform.GetPosition(entry))
		if !aabb.Intersects(region) {
			return
		}
		callback(entry)
	})
}

func (w *world) render(screen *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) []*donburi.Entry {
	w.queryResults = w.queryResults[:0]
	w.QueryRegion(viewport, func(entry *donburi.Entry) {
		if !rendering.IsVisible(entry) {
			return
		}
		w.queryResults = append(w.queryResults, entry)
	})
	w.renderer.render(w.queryResults, screen, viewport, viewMatrix)
	return w.queryResults
}
