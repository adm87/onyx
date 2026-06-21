package ecs

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/hashgrid"
	"github.com/yohamta/donburi"
)

type ECSPartitioner struct {
	entities *hashgrid.HashGrid[donburi.Entity]
	indexing map[donburi.Entity]uint64
}

func NewECSPartitioner(resolution int, padding hashgrid.Padding) *ECSPartitioner {
	return &ECSPartitioner{
		entities: hashgrid.New[donburi.Entity](resolution, padding),
		indexing: make(map[donburi.Entity]uint64),
	}
}

func (p *ECSPartitioner) Add(entity donburi.Entity, area geom.AABB) (uint64, bool) {
	if _, exists := p.indexing[entity]; exists {
		return 0, false
	}

	id := p.entities.Insert(entity, area)
	p.indexing[entity] = id

	return id, true
}

func (p *ECSPartitioner) Remove(entity donburi.Entity) {
	id, exists := p.indexing[entity]
	if !exists {
		return
	}
	p.RemoveByIndex(id)
}

func (p *ECSPartitioner) RemoveByIndex(id uint64) {
	entity, exists := p.entities.Remove(id)
	if !exists {
		return
	}
	delete(p.indexing, entity)
}

func (p *ECSPartitioner) Update(entity donburi.Entity, area geom.AABB) {
	id, exists := p.indexing[entity]
	if !exists {
		return
	}
	p.entities.Update(id, area)
}

func (p *ECSPartitioner) UpdateByIndex(id uint64, area geom.AABB) {
	p.entities.Update(id, area)
}
