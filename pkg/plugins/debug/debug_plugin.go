package debug

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

var pluginID = engine.TypeHash[DebugPlugin]()

func PluginID() uint64 {
	return pluginID
}

type DebugPlugin interface {
	engine.Plugin

	PathTransformBounds(entity donburi.Entity)

	DrawPath(target *ebiten.Image, color color.RGBA)
	ResetPath()
}

type plugin struct {
	path vector.Path
}

func NewPlugin() DebugPlugin {
	return &plugin{}
}

func (p *plugin) ID() uint64 {
	return PluginID()
}

func (p *plugin) OnRegister(game engine.Game) {
}

func (p *plugin) PathTransformBounds(entity donburi.Entity) {

}

func (p *plugin) DrawPath(target *ebiten.Image, color color.RGBA) {

}

func (p *plugin) ResetPath() {
	p.path.Reset()
}
