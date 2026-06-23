package ecs

import (
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
)

type DonburiECS struct {
	world  donburi.World
	logger engine.Logger

	renderPipeline *ECSRenderPipeline
	partitioner    *ECSPartitioner
}

func NewDonburiECS(screen engine.Screen, logger engine.Logger) *DonburiECS {
	world := donburi.NewWorld()
	partitioner := NewECSPartitioner(32, 64, 128, 256, 512)
	renderPipeline := NewECSRenderPipeline(
		world,
		screen,
		logger,
		partitioner,
	)
	return &DonburiECS{
		logger:         logger,
		world:          world,
		partitioner:    partitioner,
		renderPipeline: renderPipeline,
	}
}

func (d *DonburiECS) World() donburi.World {
	return d.world
}

func (d *DonburiECS) RenderPipeline() *ECSRenderPipeline {
	return d.renderPipeline
}

func (d *DonburiECS) Partitioner() *ECSPartitioner {
	return d.partitioner
}

func (d *DonburiECS) Add(entries ...*donburi.Entry) {
	for _, entry := range entries {
		d.AddEntry(entry)
	}
}

func (d *DonburiECS) AddEntry(entry *donburi.Entry) {
	aabb := transform.GetWorldBounds(entry)
	transform.SetIndex(entry, d.partitioner.Insert(entry.Entity(), aabb))
}

func (d *DonburiECS) Remove(entries ...*donburi.Entry) {
	for _, entry := range entries {
		d.RemoveEntry(entry)
	}
}

func (d *DonburiECS) RemoveEntry(entry *donburi.Entry) {
	d.partitioner.Remove(entry.Entity())
}

func (d *DonburiECS) Update(entries ...*donburi.Entry) {
	for _, entry := range entries {
		d.UpdateEntry(entry)
	}
}

func (d *DonburiECS) UpdateEntry(entry *donburi.Entry) {
	aabb := transform.GetWorldBounds(entry)
	index := d.partitioner.Update(entry.Entity(), aabb)
	transform.SetIndex(entry, index)
}

func (d *DonburiECS) QueryAll(area geom.AABB, fn func(entity donburi.Entity)) {
	d.partitioner.Query(area, fn)
}

func (d *DonburiECS) QueryResolution(area geom.AABB, fn func(entity donburi.Entity)) {
	partition, _ := d.partitioner.nearestPartition(area)
	partition.Query(area, fn)
}
