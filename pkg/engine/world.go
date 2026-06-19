package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/hashgrid"
	"github.com/yohamta/donburi"
)

type World interface {
	Add(entry *donburi.Entry)
	Remove(entry *donburi.Entry)
	Update(entry *donburi.Entry)
	Query(region geom.AABB, callback func(*donburi.Entry))
}

type world struct {
	ecs      donburi.World
	entities *hashgrid.HashGrid[donburi.Entity]
}

func newWorld(ecs donburi.World) *world {
	return &world{
		ecs:      ecs,
		entities: hashgrid.New[donburi.Entity](32, hashgrid.Padding{}),
	}
}

func (w *world) Add(entry *donburi.Entry) {
	if transform.GetIndexing(entry) != 0 {
		return
	}

	bounds := transform.GetWorldBounds(entry)
	id := w.entities.Insert(entry.Entity(), bounds)

	transform.SetIndexing(entry, id)
}

func (w *world) Remove(entry *donburi.Entry) {
	id := transform.GetIndexing(entry)
	if id == 0 {
		return
	}

	w.entities.Remove(id)
	transform.SetIndexing(entry, 0)
}

func (w *world) Update(entry *donburi.Entry) {
	id := transform.GetIndexing(entry)
	if id == 0 {
		return
	}

	bounds := transform.GetWorldBounds(entry)
	w.entities.Update(id, bounds)
}

func (w *world) Query(region geom.AABB, callback func(*donburi.Entry)) {
	w.entities.Query(region, func(entity donburi.Entity) {
		entry := w.ecs.Entry(entity)
		callback(entry)
	})
}

func (w *world) queryInto(region geom.AABB, result []*donburi.Entry) []*donburi.Entry {
	w.entities.Query(region, func(entity donburi.Entity) {
		entry := w.ecs.Entry(entity)
		result = append(result, entry)
	})
	return result
}
