package engine

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/quadtree"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var SceneObjectTag = donburi.NewTag()

type Scene interface {
	SyncEntities(entities []donburi.Entity)
	Render(screen *ebiten.Image, region geom.AABB, view ebiten.GeoM) error
	World() donburi.World
}

type scene struct {
	world    donburi.World
	quadtree quadtree.Quadtree
}

func newScene() Scene {
	quadtree := quadtree.New()

	world := donburi.NewWorld()
	world.OnCreate(func(world donburi.World, entity donburi.Entity) {
		onCreated(world, entity, quadtree)
	})
	world.OnRemove(func(world donburi.World, entity donburi.Entity) {
		onRemoved(world, entity, quadtree)
	})

	return &scene{
		quadtree: quadtree,
		world:    world,
	}
}

func (s *scene) World() donburi.World {
	return s.world
}

func (s *scene) Render(screen *ebiten.Image, region geom.AABB, view ebiten.GeoM) error {
	// TODO - Query quadtree for renderables.
	return nil
}

func (s *scene) SyncEntities(entities []donburi.Entity) {
	for _, entity := range entities {
		if s.world.Entry(entity).HasComponent(SceneObjectTag) {
			// TODO - Reinsert into quadtree
		}
	}
}

func onCreated(world donburi.World, entity donburi.Entity, quadtree quadtree.Quadtree) {
	if world.Entry(entity).HasComponent(SceneObjectTag) {
		// TODO - Add to quadtree
	}
}

func onRemoved(world donburi.World, entity donburi.Entity, quadtree quadtree.Quadtree) {
	if world.Entry(entity).HasComponent(SceneObjectTag) {
		// TODO - Remove from quadtree
	}
}
