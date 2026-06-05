package engine

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/spatialhash"
	"github.com/yohamta/donburi"
)

type renderableIndexing map[donburi.Entity]spatialhash.SpatialIndex

type renderables struct {
	partitioning *spatialhash.SpatialHash[donburi.Entity]
	indexing     renderableIndexing
}

func newRenderables() *renderables {
	return &renderables{
		partitioning: spatialhash.New[donburi.Entity](
			spatialhash.WithResolutions(16, 32, 64),
			spatialhash.WithCapacity(100),
		),
		indexing: make(renderableIndexing),
	}
}

func (r *renderables) add(entry *donburi.Entry, aabb geom.AABB) bool {
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

func (r *renderables) remove(entry *donburi.Entry) bool {
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

func (r *renderables) update(entry *donburi.Entry, aabb geom.AABB) bool {
	entity := entry.Entity()

	index, exists := r.indexing[entity]
	if !exists {
		return false // Entity is not indexed, cannot update
	}

	return r.partitioning.Reinsert(index, aabb)
}
