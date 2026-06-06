package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/shapes"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type World interface {
	Add(entry *donburi.Entry)
	AddMany(entries ...*donburi.Entry)

	Remove(entry *donburi.Entry)
	RemoveMany(entries ...*donburi.Entry)

	Update(entry *donburi.Entry)
	UpdateMany(entries ...*donburi.Entry)
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

func (w *world) AddMany(entries ...*donburi.Entry) {
	for i := range entries {
		w.Add(entries[i])
	}
}

func (w *world) RemoveMany(entries ...*donburi.Entry) {
	for i := range entries {
		w.Remove(entries[i])
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

func (w *world) UpdateMany(entries ...*donburi.Entry) {
	for i := range entries {
		w.Update(entries[i])
	}
}

func (w *world) render(screen *ebiten.Image, viewPort geom.AABB, viewMatrix ebiten.GeoM) error {
	return w.renderer.render(w.ecs, screen, viewPort, viewMatrix)
}
