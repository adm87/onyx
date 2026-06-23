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
		staticPartitioner:  ecs.NewECSPartitioner(32, 64, 128),
		dynamicPartitioner: ecs.NewECSPartitioner(32, 64, 128),
	}
}

func (d *CollisionWorld) Add(entries ...*donburi.Entry) {
	for _, entry := range entries {
		d.AddEntry(entry)
	}
}

func (d *CollisionWorld) AddEntry(entry *donburi.Entry) {
	if !entry.HasComponent(Collision) {
		return
	}
	var aabb geom.AABB
	if entry.HasComponent(Collision) {
		aabb = GetCollider(entry)
	} else {
		aabb = transform.GetWorldBounds(entry)
	}
	switch GetCollisionType(entry) {
	case CollisionTypeStatic:
		SetCollisionIndex(entry, d.staticPartitioner.Insert(entry.Entity(), aabb))
	case CollisionTypeDynamic:
		SetCollisionIndex(entry, d.dynamicPartitioner.Insert(entry.Entity(), aabb))
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
	var aabb geom.AABB
	if entry.HasComponent(Collision) {
		aabb = GetCollider(entry)
	} else {
		aabb = transform.GetWorldBounds(entry)
	}
	switch GetCollisionType(entry) {
	case CollisionTypeStatic:
		SetCollisionIndex(entry, d.staticPartitioner.Update(entry.Entity(), aabb))
	case CollisionTypeDynamic:
		SetCollisionIndex(entry, d.dynamicPartitioner.Update(entry.Entity(), aabb))
	}
}

func (d *CollisionWorld) QueryAll(area geom.AABB, fn func(entity donburi.Entity)) {
	d.staticPartitioner.Query(area, fn)
	d.dynamicPartitioner.Query(area, fn)
}
