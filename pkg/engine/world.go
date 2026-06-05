package engine

import "github.com/yohamta/donburi"

type World interface {
	Collision() Collision
	ECS() donburi.World
}

type world struct {
	ecs donburi.World

	collision *collision
}

func newWorld() *world {
	return &world{
		ecs:       donburi.NewWorld(),
		collision: newCollision(),
	}
}

func (w *world) ECS() donburi.World {
	return w.ecs
}

func (w *world) Collision() Collision {
	return w.collision
}
