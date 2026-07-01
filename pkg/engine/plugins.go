package engine

import (
	"fmt"
	"reflect"

	"github.com/adm87/onyx/pkg/engine/assert"
)

type Plugin interface {
	OnRegister(game Game)
	ID() uint64
}

type Plugins interface {
	GetPluginByID(pluginID uint64) (Plugin, bool)
}

type plugins struct {
	plugins map[uint64]Plugin
}

func newPlugins() *plugins {
	return &plugins{
		plugins: make(map[uint64]Plugin),
	}
}

func (p *plugins) add(plugin Plugin) {
	pluginID := plugin.ID()
	if existingPlugin, exists := p.plugins[pluginID]; exists {
		assert.Fatal(fmt.Errorf("plugin ID %d already registered to plugin of type %s", pluginID, reflect.TypeOf(existingPlugin).String()))
	}
	p.plugins[pluginID] = plugin
}

func (p *plugins) Register(game Game) {
	for _, plugin := range p.plugins {
		plugin.OnRegister(game)
	}
}

func (p *plugins) GetPluginByID(pluginID uint64) (Plugin, bool) {
	plugin, exists := p.plugins[pluginID]
	return plugin, exists
}

func GetPlugin[T Plugin](game Game, pluginID uint64) T {
	if plugin, exists := game.Plugins().GetPluginByID(pluginID); exists {
		return assert.Type[T](plugin)
	}
	panic(fmt.Errorf("plugin with ID %d not found", pluginID))
}
