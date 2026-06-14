package spatialhash

import (
	"math"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
)

func encodeCell(cellX, cellY int) uint64 {
	return uint64(uint32(cellX))<<32 | uint64(uint32(cellY))
}

type Padding struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

type SpatialHash[T comparable] struct {
	store      *slotmap.SlotMap[T]
	cellCache  []uint64
	grid       map[uint64][]uint64
	cells      map[uint64][]uint64
	querySeen  map[uint64]struct{}
	resolution int
	padding    Padding
}

func New[T comparable](resolution int, padding Padding) *SpatialHash[T] {
	return &SpatialHash[T]{
		store:      slotmap.New[T](0),
		cellCache:  make([]uint64, 0),
		grid:       make(map[uint64][]uint64),
		cells:      make(map[uint64][]uint64),
		querySeen:  make(map[uint64]struct{}),
		resolution: resolution,
		padding:    padding,
	}
}

func (sh *SpatialHash[T]) Insert(item T, area *geom.AABB) uint64 {
	id := sh.store.Insert(item)

	sh.cacheCells(area)
	for _, cell := range sh.cellCache {
		sh.grid[cell] = append(sh.grid[cell], id)
		sh.cells[id] = append(sh.cells[id], cell)
	}

	return id
}

func (sh *SpatialHash[T]) Remove(id uint64) {
	_, exists := sh.store.Get(id)
	if !exists {
		return
	}

	cells, exists := sh.cells[id]
	if !exists {
		return
	}

	for _, cell := range cells {
		ids := sh.grid[cell]
		for i, cellID := range ids {
			if cellID == id {
				sh.grid[cell] = append(sh.grid[cell][:i], sh.grid[cell][i+1:]...)
				break
			}
		}
	}

	delete(sh.cells, id)
	sh.store.Delete(id)
}

func (sh *SpatialHash[T]) Update(id uint64, area *geom.AABB) {
	_, exists := sh.store.Get(id)
	if !exists {
		return
	}

	for _, cell := range sh.cells[id] {
		ids := sh.grid[cell]
		for i, cellID := range ids {
			if cellID == id {
				sh.grid[cell] = append(sh.grid[cell][:i], sh.grid[cell][i+1:]...)
				break
			}
		}
	}

	sh.cacheCells(area)
	sh.cells[id] = sh.cells[id][:0]

	for _, cell := range sh.cellCache {
		sh.grid[cell] = append(sh.grid[cell], id)
		sh.cells[id] = append(sh.cells[id], cell)
	}
}

func (sh *SpatialHash[T]) Query(area *geom.AABB, fn func(item T)) {
	sh.cacheCells(area)
	clear(sh.querySeen)

	for _, cell := range sh.cellCache {
		ids, exists := sh.grid[cell]
		if !exists {
			continue
		}

		for _, id := range ids {
			if _, alreadySeen := sh.querySeen[id]; alreadySeen {
				continue
			}
			sh.querySeen[id] = struct{}{}

			item, exists := sh.store.Get(id)
			if !exists {
				continue
			}

			fn(item)
		}
	}
}

func (sh *SpatialHash[T]) cacheCells(area *geom.AABB) {
	sh.cellCache = sh.cellCache[:0]

	minX := math.Floor(area.Min.X/float64(sh.resolution)) - float64(sh.padding.Left)
	minY := math.Floor(area.Min.Y/float64(sh.resolution)) - float64(sh.padding.Top)
	maxX := math.Floor(area.Max.X/float64(sh.resolution)) + float64(sh.padding.Right)
	maxY := math.Floor(area.Max.Y/float64(sh.resolution)) + float64(sh.padding.Bottom)

	for x := int(minX); x <= int(maxX); x++ {
		for y := int(minY); y <= int(maxY); y++ {
			sh.cellCache = append(sh.cellCache, encodeCell(x, y))
		}
	}
}
