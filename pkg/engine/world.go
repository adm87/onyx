package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/spatialhash"
	"github.com/yohamta/donburi"
)

type World interface {
	Add(entry *donburi.Entry)
	Remove(entry *donburi.Entry)
	Update(entry *donburi.Entry)

	QueryInto(ecs donburi.World, region geom.AABB, results []donburi.Entity) []donburi.Entity
	QueryRegion(ecs donburi.World, region geom.AABB, callback func(entry *donburi.Entry))
}

type world struct {
	entities     *spatialhash.SpatialHash[donburi.Entity]
	queryResults []*donburi.Entry
}

var worldIndexing = donburi.NewComponentType[uint64]()

func newWorld() *world {
	return &world{
		entities:     spatialhash.New[donburi.Entity](16, spatialhash.Padding{}),
		queryResults: make([]*donburi.Entry, 0, 100),
	}
}

func (w *world) Add(entry *donburi.Entry) {
	aabb := transform.GetWorldBounds(entry)
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
	aabb := transform.GetWorldBounds(entry)
	w.entities.Update(*index, aabb)
}

func (w *world) UpdateBounds(entry *donburi.Entry, bounds geom.AABB) {
	index := worldIndexing.Get(entry)
	w.entities.Update(*index, bounds)
}

func (w *world) QueryInto(ecs donburi.World, region geom.AABB, results []donburi.Entity) []donburi.Entity {
	return w.entities.QueryInto(region, results)
}

func (w *world) QueryRegion(ecs donburi.World, region geom.AABB, callback func(entry *donburi.Entry)) {
	w.entities.Query(region, func(entity donburi.Entity) {
		callback(ecs.Entry(entity))
	})
}
