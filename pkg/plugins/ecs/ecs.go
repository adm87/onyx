package ecs

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning/hashgrid"
	"github.com/adm87/onyx/pkg/plugins/ecs/image"
	"github.com/adm87/onyx/pkg/plugins/ecs/tiled"
	"github.com/adm87/onyx/pkg/plugins/ecs/transform"
	"github.com/yohamta/donburi"

	imageplugin "github.com/adm87/onyx/pkg/plugins/images"
	tiledplugin "github.com/adm87/onyx/pkg/plugins/tiled"
)

type DonburiECSPlugin struct {
	world  donburi.World
	logger engine.Logger

	factory        *ECSFactory
	renderPipeline *ECSRenderPipeline
	partitioner    *ECSPartitioner
}

func NewDonburiECSPlugin(
	screen engine.Screen,
	logger engine.Logger,
	imageAssets *imageplugin.ImageAssets,
	tiledAssets *tiledplugin.TiledAssets) *DonburiECSPlugin {

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
		factory: &ECSFactory{
			imageAssets: imageAssets,
			tiledAssets: tiledAssets,
			imageRendererType: renderPipeline.AddAdapter(image.NewImageRenderer(
				imageAssets,
			)),
			tiledRendererType: renderPipeline.AddAdapter(tiled.NewTiledRenderer(
				imageAssets,
				tiledAssets,
			)),
			partitioner: partitioner,
		},
	}
}

func (d *DonburiECSPlugin) World() donburi.World {
	return d.world
}

func (d *DonburiECSPlugin) RenderPipeline() *ECSRenderPipeline {
	return d.renderPipeline
}

func (d *DonburiECSPlugin) Factory() *ECSFactory {
	return d.factory
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
