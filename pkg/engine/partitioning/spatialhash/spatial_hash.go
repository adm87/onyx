package spatialhash

import (
	"math"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap.go"
)

type (
	SpatialIndex uint64
	spatialCoord uint64
)

func encodeCoord(cellX, cellY int64) spatialCoord {
	x := uint64(cellX) & 0xFFFFFFFF
	y := uint64(cellY) & 0xFFFFFFFF
	return spatialCoord(x | (y << 32))
}

type spatialEntry struct {
	key   slotmap.SlotKey
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

	storage   *slotmap.SlotMap[T]
	nextIndex SpatialIndex

	cells []spatialCoord
	seen  map[SpatialIndex]struct{}

	padding [4]int
}

func New[T comparable](opts ...Option) *SpatialHash[T] {
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	grids := make([]spatialGrid, len(options.Resolutions))
	for i, res := range options.Resolutions {
		grids[i] = spatialGrid{
			cellSize: res,
			cells:    make(map[spatialCoord][]SpatialIndex),
		}
	}
	return &SpatialHash[T]{
		grids:   grids,
		index:   make(map[SpatialIndex]spatialEntry, options.Capacity),
		seen:    make(map[SpatialIndex]struct{}, options.Capacity),
		storage: slotmap.New[T](options.Capacity),
		padding: options.Padding,
	}
}

func (h *SpatialHash[T]) Insert(value T, aabb geom.AABB) (SpatialIndex, bool) {
	key, ok := h.storage.Insert(value)
	if !ok {
		return 0, false
	}

	index := h.nextIndex
	h.nextIndex++

	entry := spatialEntry{key: key}
	h.addToGrid(index, &entry, aabb)

	h.index[index] = entry
	return index, true
}

func (h *SpatialHash[T]) Remove(index SpatialIndex) bool {
	entry, exists := h.index[index]
	if !exists {
		return false
	}

	h.removeFromGrid(index, &entry)
	delete(h.index, index)
	return h.storage.Remove(entry.key)
}

func (h *SpatialHash[T]) Reinsert(index SpatialIndex, aabb geom.AABB) bool {
	entry, exists := h.index[index]
	if !exists {
		return false
	}

	h.removeFromGrid(index, &entry)
	h.addToGrid(index, &entry, aabb)
	h.index[index] = entry
	return true
}

func (h *SpatialHash[T]) addToGrid(index SpatialIndex, entry *spatialEntry, aabb geom.AABB) {
	width := aabb.Max.X - aabb.Min.X
	height := aabb.Max.Y - aabb.Min.Y

	idx := h.getNearestGrid(max(width, height))
	grid := &h.grids[idx]
	h.getCells(aabb, grid)

	entry.grid = idx
	entry.cells = entry.cells[:0]
	for _, cell := range h.cells {
		grid.cells[cell] = append(grid.cells[cell], index)
		entry.cells = append(entry.cells, cell)
	}
}

func (h *SpatialHash[T]) removeFromGrid(index SpatialIndex, entry *spatialEntry) {
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
	cellMinX := math.Floor(aabb.Min.X/grid.cellSize) - float64(h.padding[0])
	cellMinY := math.Floor(aabb.Min.Y/grid.cellSize) - float64(h.padding[1])
	cellMaxX := math.Floor(aabb.Max.X/grid.cellSize) + float64(h.padding[2])
	cellMaxY := math.Floor(aabb.Max.Y/grid.cellSize) + float64(h.padding[3])

	h.cells = h.cells[:0]
	for x := cellMinX; x <= cellMaxX; x++ {
		for y := cellMinY; y <= cellMaxY; y++ {
			h.cells = append(h.cells, encodeCoord(int64(x), int64(y)))
		}
	}
	return h.cells
}
