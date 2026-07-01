package collision

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/ecs"
	"github.com/yohamta/donburi"
)

var pluginID = engine.TypeHash[CollisionPlugin]()

func PluginID() uint64 {
	return pluginID
}

type CollisionPlugin interface {
	engine.Plugin
}

type plugin struct {
}

func NewPlugin() CollisionPlugin {
	return &plugin{}
}

func (p *plugin) OnRegister(game engine.Game) {
	ecsPlugin := engine.GetPlugin[ecs.ECSPlugin](game, ecs.PluginID())
	ecsPlugin.AddECSCallbacks(ecs.ECSCallbacks{
		Added:   p.addCollisionEntries,
		Removed: p.removeCollisionEntries,
		Updated: p.updateCollisionEntries,
	})
}

func (p *plugin) ID() uint64 {
	return PluginID()
}

func (p *plugin) addCollisionEntries(entries []*donburi.Entry) {

}

func (p *plugin) removeCollisionEntries(entries []*donburi.Entry) {

}

func (p *plugin) updateCollisionEntries(entries []*donburi.Entry) {

}
