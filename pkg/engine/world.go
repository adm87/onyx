package engine

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type World interface {
	ECS() donburi.World

	Add(entry *donburi.Entry)
	Remove(entry *donburi.Entry)
	Update(entry *donburi.Entry)
}

type world struct {
	ecs donburi.World

	collision *collision
	renderer  *renderer
}

func newWorld(collision *collision, renderer *renderer) *world {
	ecs := donburi.NewWorld()
	return &world{
		ecs:       ecs,
		collision: collision,
		renderer:  renderer,
	}
}

func (w *world) ECS() donburi.World {
	return w.ecs
}

func (w *world) Add(entry *donburi.Entry) {

}

func (w *world) Remove(entry *donburi.Entry) {

}

func (w *world) Update(entry *donburi.Entry) {

}

func (w *world) render(screen *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) {

}
