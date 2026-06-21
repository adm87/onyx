package ecs

import (
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/hashgrid"
	"github.com/yohamta/donburi"
)

type DonburiECSPlugin struct {
	world  donburi.World
	logger engine.Logger

	renderPipeline *ECSRenderPipeline
	partitioner    *ECSPartitioner
}

func NewDonburiECSPlugin(screen engine.Screen, logger engine.Logger) *DonburiECSPlugin {
	world := donburi.NewWorld()
	partitioner := NewECSPartitioner(
		256,
		hashgrid.Padding{},
	)
	renderPipeline := NewECSRenderPipeline(
		world,
		screen,
		logger,
		partitioner,
	)

	return &DonburiECSPlugin{
		logger:         logger,
		world:          world,
		partitioner:    partitioner,
		renderPipeline: renderPipeline,
	}
}

func (d *DonburiECSPlugin) World() donburi.World {
	return d.world
}

func (d *DonburiECSPlugin) RenderPipeline() *ECSRenderPipeline {
	return d.renderPipeline
}

func (d *DonburiECSPlugin) Add(entries ...*donburi.Entry) {
	for i := range entries {
		aabb := transform.GetWorldBounds(entries[i])

		index, ok := d.partitioner.Add(entries[i].Entity(), aabb)
		if !ok {
			d.logger.Warn("Entity already exists in partitioner, skipping addition.")
			continue
		}

		transform.SetIndex(entries[i], index)
	}
}

func (d *DonburiECSPlugin) Update(entries ...*donburi.Entry) {
	for i := range entries {
		aabb := transform.GetWorldBounds(entries[i])
		d.partitioner.Update(entries[i].Entity(), aabb)
	}
}

func (d *DonburiECSPlugin) Remove(entries ...*donburi.Entry) {
	for i := range entries {
		entity := entries[i].Entity()
		d.partitioner.Remove(entity)
		d.world.Remove(entity)
	}
}

func (d *DonburiECSPlugin) Query(area geom.AABB, callback func(entity donburi.Entity)) {
	d.partitioner.entities.Query(area, callback)
}
