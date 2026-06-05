package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/yohamta/donburi"
)

type World interface {
	Add(entry *donburi.Entry)
	Remove(entry *donburi.Entry)
	Update(entry *donburi.Entry)

	Collision() Collision
	ECS() donburi.World
}

type world struct {
	ecs donburi.World

	collision   *collision
	renderables *renderables
}

func newWorld() *world {
	return &world{
		ecs:         donburi.NewWorld(),
		collision:   newCollision(),
		renderables: newRenderables(),
	}
}

func (w *world) ECS() donburi.World {
	return w.ecs
}

func (w *world) Collision() Collision {
	return w.collision
}

func (w *world) Add(entry *donburi.Entry) {
	position := transform.GetPosition(entry)
	aabb := colliders.GetAABB(entry).Translate(position)

	if entry.HasComponent(colliders.Collision) {
		w.collision.add(entry, aabb)
	}
	if entry.HasComponent(rendering.Renderer) {
		w.renderables.add(entry, aabb)
	}
}

func (w *world) Remove(entry *donburi.Entry) {
	w.collision.remove(entry)
	w.renderables.remove(entry)
}

func (w *world) Update(entry *donburi.Entry) {
	position := transform.GetPosition(entry)
	aabb := colliders.GetAABB(entry).Translate(position)

	w.collision.update(entry, aabb)
	w.renderables.update(entry, aabb)
}
