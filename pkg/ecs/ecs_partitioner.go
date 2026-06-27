package ecs

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/hashgrid"
	"github.com/yohamta/donburi"
)

type partitionIndex struct {
	idx       uint64
	partition int
}

type ECSPartitioner struct {
	partitions []*hashgrid.HashGrid[donburi.Entity]
	indexing   map[donburi.Entity]partitionIndex

	queryGen  uint32
	querySeen map[donburi.Entity]uint32
}

func NewECSPartitioner(resolutions ...int) *ECSPartitioner {
	partitions := make([]*hashgrid.HashGrid[donburi.Entity], len(resolutions))
	for i, res := range resolutions {
		partitions[i] = hashgrid.New[donburi.Entity](res)
	}
	return &ECSPartitioner{
		partitions: partitions,
		indexing:   make(map[donburi.Entity]partitionIndex),
		querySeen:  make(map[donburi.Entity]uint32),
	}
}

func (p *ECSPartitioner) GetPartition(i int) *hashgrid.HashGrid[donburi.Entity] {
	if i < 0 || i >= len(p.partitions) {
		return nil
	}
	return p.partitions[i]
}

func (p *ECSPartitioner) Insert(entity donburi.Entity, area geom.AABB) uint64 {
	if index, exists := p.indexing[entity]; exists {
		return index.idx
	}
	partition, i := p.nearestPartition(area)

	id := partition.Insert(entity, area)
	p.indexing[entity] = partitionIndex{
		idx:       id,
		partition: i,
	}

	return id
}

func (p *ECSPartitioner) Remove(entity donburi.Entity) {
	index, exists := p.indexing[entity]
	if !exists {
		return
	}

	partition := p.partitions[index.partition]
	partition.Remove(index.idx)

	delete(p.indexing, entity)
}

func (p *ECSPartitioner) Update(entity donburi.Entity, area geom.AABB) uint64 {
	index, exists := p.indexing[entity]
	if !exists {
		return p.Insert(entity, area)
	}

	partition, i := p.nearestPartition(area)
	if i == index.partition {
		partition.Update(index.idx, area)
		return index.idx
	}

	oldPartition := p.partitions[index.partition]
	oldPartition.Remove(index.idx)

	id := partition.Insert(entity, area)
	p.indexing[entity] = partitionIndex{
		idx:       id,
		partition: i,
	}

	return id
}

func (p *ECSPartitioner) Query(area geom.AABB, callback func(donburi.Entity)) {
	p.queryGen++
	for _, partition := range p.partitions {
		partition.Query(area, func(entity donburi.Entity) {
			if p.querySeen[entity] == p.queryGen {
				return
			}
			p.querySeen[entity] = p.queryGen
			callback(entity)
		})
	}
}

func (p *ECSPartitioner) nearestPartition(aabb geom.AABB) (*hashgrid.HashGrid[donburi.Entity], int) {
	resolution := int(max(aabb.Width(), aabb.Height()))
	for i, partition := range p.partitions {
		if resolution <= partition.Resolution() {
			return partition, i
		}
	}
	i := len(p.partitions) - 1
	return p.partitions[i], i
}
