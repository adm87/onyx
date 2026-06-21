package hashgrid

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

type HashGrid[T comparable] struct {
	store      *slotmap.SlotMap[T]
	queryGen   uint32
	resolution int
	padding    Padding
	cellCache  []uint64
	grid       map[uint64][]uint64
	cells      map[uint64][]uint64
	querySeen  map[uint64]uint32
}

func New[T comparable](resolution int, padding Padding) *HashGrid[T] {
	return &HashGrid[T]{
		store:      slotmap.New[T](0),
		cellCache:  make([]uint64, 0),
		grid:       make(map[uint64][]uint64),
		cells:      make(map[uint64][]uint64),
		querySeen:  make(map[uint64]uint32),
		resolution: resolution,
		padding:    padding,
	}
}

func (sh *HashGrid[T]) Resolution() int {
	return sh.resolution
}

func (sh *HashGrid[T]) Insert(item T, area geom.AABB) uint64 {
	id := sh.store.Insert(item)

	sh.cacheCells(area)
	for _, cell := range sh.cellCache {
		sh.grid[cell] = append(sh.grid[cell], id)
		sh.cells[id] = append(sh.cells[id], cell)
	}

	return id
}

func (sh *HashGrid[T]) Remove(id uint64) (T, bool) {
	item, exists := sh.store.Get(id)
	if !exists {
		var zero T
		return zero, false
	}

	cells, exists := sh.cells[id]
	if !exists {
		var zero T
		return zero, false
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

	return item, true
}

func (sh *HashGrid[T]) Update(id uint64, area geom.AABB) {
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

func (sh *HashGrid[T]) Query(area geom.AABB, fn func(item T)) {
	sh.cacheCells(area)
	sh.queryGen++
	for _, cell := range sh.cellCache {
		ids, exists := sh.grid[cell]
		if !exists {
			continue
		}

		for _, id := range ids {
			if sh.querySeen[id] == sh.queryGen {
				continue
			}
			sh.querySeen[id] = sh.queryGen

			item, exists := sh.store.Get(id)
			if !exists {
				continue
			}

			fn(item)
		}
	}
}

func (sh *HashGrid[T]) GetCellsWithin(area geom.AABB) []geom.AABB {
	sh.cacheCells(area)
	cells := make([]geom.AABB, 0, len(sh.cellCache))
	for _, cell := range sh.cellCache {
		cellX := int(int32(cell >> 32))
		cellY := int(int32(cell & 0xFFFFFFFF))

		cells = append(cells, geom.AABB{
			Min: geom.Vec2{
				X: float64(cellX * sh.resolution),
				Y: float64(cellY * sh.resolution),
			},
			Max: geom.Vec2{
				X: float64((cellX + 1) * sh.resolution),
				Y: float64((cellY + 1) * sh.resolution),
			},
		})
	}
	return cells
}

func (sh *HashGrid[T]) cacheCells(area geom.AABB) {
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
