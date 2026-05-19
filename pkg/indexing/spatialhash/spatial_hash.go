package spatialhash

import "github.com/adm87/onyx/pkg/engine/geom"

type SpatialHash[T comparable] struct {
}

func New[T comparable]() *SpatialHash[T] {
	return &SpatialHash[T]{}
}

func (sh *SpatialHash[T]) Insert(item T, aabb geom.AABB) bool {
	return false
}

func (sh *SpatialHash[T]) Remove(item T) bool {
	return false
}

func (sh *SpatialHash[T]) Query(region geom.AABB) []T {
	return nil
}

func (sh *SpatialHash[T]) BuildIndexing() {

}
