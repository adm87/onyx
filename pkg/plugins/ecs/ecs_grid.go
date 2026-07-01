package ecs

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/hashgrid"
	"github.com/yohamta/donburi"
)

type gridIndex struct {
	idx  uint64
	grid int
}

type ECSGrid struct {
	grid     []*hashgrid.HashGrid[donburi.Entity]
	indexing map[donburi.Entity]gridIndex

	queryGen  uint32
	querySeen map[donburi.Entity]uint32
}

func NewEntityGrid(resolutions ...int) *ECSGrid {
	grids := make([]*hashgrid.HashGrid[donburi.Entity], len(resolutions))
	for i, res := range resolutions {
		grids[i] = hashgrid.New[donburi.Entity](res)
	}
	return &ECSGrid{
		grid:      grids,
		indexing:  make(map[donburi.Entity]gridIndex),
		querySeen: make(map[donburi.Entity]uint32),
	}
}

func (p *ECSGrid) GetGrid(i int) *hashgrid.HashGrid[donburi.Entity] {
	if i < 0 || i >= len(p.grid) {
		return nil
	}
	return p.grid[i]
}

func (p *ECSGrid) Insert(entity donburi.Entity, area geom.AABB) uint64 {
	if index, exists := p.indexing[entity]; exists {
		return index.idx
	}
	grid, i := p.nearestGrid(area)

	id := grid.Insert(entity, area)
	p.indexing[entity] = gridIndex{
		idx:  id,
		grid: i,
	}

	return id
}

func (p *ECSGrid) Remove(entity donburi.Entity) {
	index, exists := p.indexing[entity]
	if !exists {
		return
	}

	grid := p.grid[index.grid]
	grid.Remove(index.idx)

	delete(p.indexing, entity)
}

func (p *ECSGrid) Update(entity donburi.Entity, area geom.AABB) uint64 {
	index, exists := p.indexing[entity]
	if !exists {
		return p.Insert(entity, area)
	}

	grid, i := p.nearestGrid(area)
	if i == index.grid {
		grid.Update(index.idx, area)
		return index.idx
	}

	oldGrid := p.grid[index.grid]
	oldGrid.Remove(index.idx)

	id := grid.Insert(entity, area)
	p.indexing[entity] = gridIndex{
		idx:  id,
		grid: i,
	}

	return id
}

func (p *ECSGrid) Query(area geom.AABB, callback func(donburi.Entity)) {
	p.queryGen++
	for _, grid := range p.grid {
		grid.Query(area, func(entity donburi.Entity) {
			if p.querySeen[entity] == p.queryGen {
				return
			}
			p.querySeen[entity] = p.queryGen
			callback(entity)
		})
	}
}

func (p *ECSGrid) nearestGrid(aabb geom.AABB) (*hashgrid.HashGrid[donburi.Entity], int) {
	resolution := int(max(aabb.Width(), aabb.Height()))
	for i, grid := range p.grid {
		if resolution <= grid.Resolution() {
			return grid, i
		}
	}
	i := len(p.grid) - 1
	return p.grid[i], i
}
