package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/shapes"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type World interface {
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
	return &world{
		ecs:       donburi.NewWorld(),
		collision: collision,
		renderer:  renderer,
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
	aabb := shapes.GetAABB(entry).Translate(position)

	if entry.HasComponent(colliders.Collision) {
		w.collision.add(entry, aabb)
	}
	if entry.HasComponent(rendering.Renderer) {
		w.renderer.addRenderable(entry, aabb)
	}
}

func (w *world) Remove(entry *donburi.Entry) {
	w.collision.remove(entry)
	w.renderer.removeRenderable(entry)
}

func (w *world) Update(entry *donburi.Entry) {
	position := transform.GetPosition(entry)
	aabb := shapes.GetAABB(entry).Translate(position)

	w.collision.update(entry, aabb)
	w.renderer.updateRenderable(entry, aabb)
}

func (w *world) render(screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
	return w.renderer.render(w.ecs, screen, viewMatrix)
}
