package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Scene interface {
	Render(screen *ebiten.Image) error
	World() donburi.World
}

type scene struct {
	world donburi.World
}

func newScene() Scene {
	s := &scene{
		world: donburi.NewWorld(),
	}
	s.world.OnCreate(s.onCreated)
	s.world.OnRemove(s.onRemoved)
	return s
}

func (s *scene) World() donburi.World {
	return s.world
}

func (s *scene) Render(screen *ebiten.Image) error {
	return nil
}

func (s *scene) onCreated(world donburi.World, entity donburi.Entity) {

}

func (s *scene) onRemoved(world donburi.World, entity donburi.Entity) {

}
