package collision

import (
	"github.com/adm87/onyx/pkg/ecs"
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
)

type CollisionWorld struct {
	staticPartitioner  *ecs.ECSPartitioner
	dynamicPartitioner *ecs.ECSPartitioner
}

func NewCollisionWorld() *CollisionWorld {
	return &CollisionWorld{
		staticPartitioner:  ecs.NewECSPartitioner(32),
		dynamicPartitioner: ecs.NewECSPartitioner(32),
	}
}

func (d *CollisionWorld) Add(entries ...*donburi.Entry) {
	for _, entry := range entries {
		d.AddEntry(entry)
	}
}

func (d *CollisionWorld) Static() *ecs.ECSPartitioner {
	return d.staticPartitioner
}

func (d *CollisionWorld) Dynamic() *ecs.ECSPartitioner {
	return d.dynamicPartitioner
}

func (d *CollisionWorld) AddEntry(entry *donburi.Entry) {
	if !entry.HasComponent(Collision) {
		return
	}
	switch GetCollisionType(entry) {
	case CollisionTypeStatic:
		d.partitionEntry(entry, d.staticPartitioner.Insert)
	case CollisionTypeDynamic:
		d.partitionEntry(entry, d.dynamicPartitioner.Insert)
	}
}

func (d *CollisionWorld) Remove(entries ...*donburi.Entry) {
	for _, entry := range entries {
		d.RemoveEntry(entry)
	}
}

func (d *CollisionWorld) RemoveEntry(entry *donburi.Entry) {
	if !entry.HasComponent(Collision) {
		return
	}
	switch GetCollisionType(entry) {
	case CollisionTypeStatic:
		d.staticPartitioner.Remove(entry.Entity())
	case CollisionTypeDynamic:
		d.dynamicPartitioner.Remove(entry.Entity())
	}
}

func (d *CollisionWorld) Update(entries ...*donburi.Entry) {
	for _, entry := range entries {
		d.UpdateEntry(entry)
	}
}

func (d *CollisionWorld) UpdateEntry(entry *donburi.Entry) {
	if !entry.HasComponent(Collision) {
		return
	}
	switch GetCollisionType(entry) {
	case CollisionTypeStatic:
		d.partitionEntry(entry, d.staticPartitioner.Update)
	case CollisionTypeDynamic:
		d.partitionEntry(entry, d.dynamicPartitioner.Update)
	}
}

func (d *CollisionWorld) QueryAll(area geom.AABB, fn func(entity donburi.Entity)) {
	d.staticPartitioner.Query(area, fn)
	d.dynamicPartitioner.Query(area, fn)
}

func (d *CollisionWorld) partitionEntry(entry *donburi.Entry, partitionFn func(entity donburi.Entity, aabb geom.AABB) uint64) {
	var aabb geom.AABB
	if entry.HasComponent(Collision) {
		aabb = GetWorldCollider(entry)
	} else {
		aabb = transform.GetWorldBounds(entry)
	}
	SetCollisionIndex(entry, partitionFn(entry.Entity(), aabb))
}
