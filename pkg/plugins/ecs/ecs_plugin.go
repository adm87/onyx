package ecs

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/ecs/transform"
	"github.com/yohamta/donburi"
)

var pluginID = engine.TypeHash[ECSPlugin]()

func PluginID() uint64 {
	return pluginID
}

type ECSCallbacks struct {
	Added   func(entries []*donburi.Entry)
	Removed func(entries []*donburi.Entry)
	Updated func(entries []*donburi.Entry)
}

type ECSPlugin interface {
	engine.Plugin

	AddECSCallbacks(callbacks ECSCallbacks)

	Add(entries ...*donburi.Entry)
	Remove(entries ...*donburi.Entry)
	Update(entries ...*donburi.Entry)

	QueryAll(area geom.AABB, fn func(entry *donburi.Entry))
	QueryResolution(area geom.AABB, fn func(entry *donburi.Entry))

	World() donburi.World
	RenderPipeline() *ECSRenderPipeline
}

type plugin struct {
	world donburi.World

	renderPipeline *ECSRenderPipeline
	grid           *ECSGrid

	callbacks []ECSCallbacks
}

func NewPlugin() ECSPlugin {
	world := donburi.NewWorld()
	ecsGrid := NewEntityGrid(32, 64, 128, 256, 512)
	renderPipeline := NewECSRenderPipeline(world, ecsGrid)
	return &plugin{
		world:          world,
		renderPipeline: renderPipeline,
		grid:           ecsGrid,
		callbacks:      make([]ECSCallbacks, 0),
	}
}

func (p *plugin) OnRegister(game engine.Game) {
	p.AddECSCallbacks(ECSCallbacks{
		Added:   p.addEntries,
		Removed: p.removeEntries,
		Updated: p.updateEntries,
	})
}

func (p *plugin) ID() uint64 {
	return PluginID()
}

func (p *plugin) RenderPipeline() *ECSRenderPipeline {
	return p.renderPipeline
}

func (p *plugin) World() donburi.World {
	return p.world
}

func (p *plugin) AddECSCallbacks(callbacks ECSCallbacks) {
	p.callbacks = append(p.callbacks, callbacks)
}

func (p *plugin) Add(entries ...*donburi.Entry) {
	for i := range p.callbacks {
		callbacks := p.callbacks[i]
		if callbacks.Added != nil {
			callbacks.Added(entries)
		}
	}
}

func (p *plugin) Remove(entries ...*donburi.Entry) {
	for i := range p.callbacks {
		callbacks := p.callbacks[i]
		if callbacks.Removed != nil {
			callbacks.Removed(entries)
		}
	}
}

func (p *plugin) Update(entries ...*donburi.Entry) {
	for i := range p.callbacks {
		callbacks := p.callbacks[i]
		if callbacks.Updated != nil {
			callbacks.Updated(entries)
		}
	}
}

func (p *plugin) addEntries(entries []*donburi.Entry) {
	for i := range entries {
		entry := entries[i]
		transform.SetIndex(entry, p.grid.Insert(entry.Entity(), transform.GetWorldBounds(entry)))
	}
}

func (p *plugin) removeEntries(entries []*donburi.Entry) {
	for i := range entries {
		entity := entries[i].Entity()
		p.grid.Remove(entity)
		p.world.Remove(entity)
	}
}

func (p *plugin) updateEntries(entries []*donburi.Entry) {
	for i := range entries {
		entry := entries[i]
		transform.SetIndex(entry, p.grid.Update(entry.Entity(), transform.GetWorldBounds(entry)))
	}
}

func (p *plugin) QueryAll(area geom.AABB, fn func(entry *donburi.Entry)) {
	p.grid.Query(area, func(entity donburi.Entity) {
		entry := p.world.Entry(entity)
		fn(entry)
	})
}

func (p *plugin) QueryResolution(area geom.AABB, fn func(entry *donburi.Entry)) {
	partition, _ := p.grid.nearestGrid(area)
	partition.Query(area, func(entity donburi.Entity) {
		entry := p.world.Entry(entity)
		fn(entry)
	})
}
