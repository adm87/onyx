package partitioning

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/storage"
)

type SpatialIndex uint64

type spatialCoord uint64

func encodeCoord(cellX, cellY int64) spatialCoord {
	x := uint64(cellX) & 0xFFFFFFFF
	y := uint64(cellY) & 0xFFFFFFFF
	return spatialCoord(x | (y << 32))
}

func decodeCoord(coord spatialCoord) (cellX, cellY int64) {
	x := int64(uint64(coord) & 0xFFFFFFFF)
	y := int64((uint64(coord) >> 32) & 0xFFFFFFFF)
	return x, y
}

type spatialEntry struct {
	key   storage.SlotKey
	cells []spatialCoord
	grid  int
}

type spatialGrid struct {
	cellSize float64
	cells    map[spatialCoord][]SpatialIndex
}

type SpatialHash[T comparable] struct {
	grids []spatialGrid
	index map[SpatialIndex]spatialEntry

	storage   *storage.SlotMap[T]
	nextIndex SpatialIndex

	cells []spatialCoord
	seen  map[SpatialIndex]struct{}
}

func NewSpatialHash[T comparable](capacity int, resolutions ...float64) *SpatialHash[T] {
	grids := make([]spatialGrid, len(resolutions))
	for i, res := range resolutions {
		grids[i] = spatialGrid{
			cellSize: res,
			cells:    make(map[spatialCoord][]SpatialIndex),
		}
	}
	return &SpatialHash[T]{
		grids:   grids,
		storage: storage.NewSlotMap[T](capacity),
		index:   make(map[SpatialIndex]spatialEntry),
		seen:    make(map[SpatialIndex]struct{}),
		cells:   make([]spatialCoord, 0),
	}
}

func (h *SpatialHash[T]) Insert(value T, aabb geom.AABB) (SpatialIndex, bool) {
	key, ok := h.storage.Insert(value)
	if !ok {
		return 0, false
	}

	width := aabb.Max.X - aabb.Min.X
	height := aabb.Max.Y - aabb.Min.Y

	idx := h.getNearestGrid(max(width, height))

	grid := &h.grids[idx]
	h.getCells(aabb, grid)

	index := h.nextIndex
	h.nextIndex++

	entry := spatialEntry{
		key:   key,
		grid:  idx,
		cells: make([]spatialCoord, len(h.cells)),
	}

	for i, cell := range h.cells {
		grid.cells[cell] = append(grid.cells[cell], index)
		entry.cells[i] = cell
	}

	h.index[index] = entry
	return index, true
}

func (h *SpatialHash[T]) Remove(index SpatialIndex) bool {
	entry, exists := h.index[index]
	if !exists {
		return false
	}

	grid := &h.grids[entry.grid]
	for _, cell := range entry.cells {
		indices := grid.cells[cell]
		for i, idx := range indices {
			if idx == index {
				grid.cells[cell] = append(indices[:i], indices[i+1:]...)
				break
			}
		}
		if len(grid.cells[cell]) == 0 {
			delete(grid.cells, cell)
		}
	}

	delete(h.index, index)
	return h.storage.Remove(entry.key)
}

// QueryNearest returns all values on the same grid resolution as the provided AABB.
// The query will stop if the provided function returns false.
func (h *SpatialHash[T]) QueryNearest(aabb geom.AABB, fn func(T) bool) {
	width := aabb.Max.X - aabb.Min.X
	height := aabb.Max.Y - aabb.Min.Y

	idx := h.getNearestGrid(max(width, height))
	grid := &h.grids[idx]

	h.getCells(aabb, grid)

	clear(h.seen)
	for _, cell := range h.cells {
		for _, index := range grid.cells[cell] {
			if _, ok := h.seen[index]; !ok {
				h.seen[index] = struct{}{}
				entry := h.index[index]
				value, ok := h.storage.Get(entry.key)
				if !ok {
					continue
				}
				if !fn(value) {
					return
				}
			}
		}
	}
}

// QueryAll returns all values on all grid resolutions that intersect with the provided AABB.
// The query will stop if the provided function returns false.
func (h *SpatialHash[T]) QueryAll(aabb geom.AABB, fn func(T) bool) {
	clear(h.seen)
	for i := range h.grids {
		grid := &h.grids[i]
		h.getCells(aabb, grid)

		for _, cell := range h.cells {
			for _, index := range grid.cells[cell] {
				if _, ok := h.seen[index]; !ok {
					h.seen[index] = struct{}{}
					entry := h.index[index]
					value, ok := h.storage.Get(entry.key)
					if !ok {
						continue
					}
					if !fn(value) {
						return
					}
				}
			}
		}
	}
}

func (h *SpatialHash[T]) getNearestGrid(size float64) int {
	for i := range h.grids {
		grid := &h.grids[i]
		if size <= grid.cellSize {
			return i
		}
	}
	return len(h.grids) - 1
}

func (h *SpatialHash[T]) getCells(aabb geom.AABB, grid *spatialGrid) []spatialCoord {
	cellMinX := int64(aabb.Min.X / grid.cellSize)
	cellMinY := int64(aabb.Min.Y / grid.cellSize)
	cellMaxX := int64(aabb.Max.X / grid.cellSize)
	cellMaxY := int64(aabb.Max.Y / grid.cellSize)

	h.cells = h.cells[:0]
	for x := cellMinX; x <= cellMaxX; x++ {
		for y := cellMinY; y <= cellMaxY; y++ {
			h.cells = append(h.cells, encodeCoord(x, y))
		}
	}
	return h.cells
}
