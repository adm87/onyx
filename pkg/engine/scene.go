package engine

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SceneIndexer interface {
	Insert(entity donburi.Entity, aabb geom.AABB) bool
	Remove(entity donburi.Entity) bool
	Query(region geom.AABB) []donburi.Entity
	BuildIndexing()
}

type Scene interface {
	Render(screen *ebiten.Image, region geom.AABB, view ebiten.GeoM) error
	World() donburi.World
}

type scene struct {
	world   donburi.World
	indexer SceneIndexer
}

func newScene(sceneIndexer SceneIndexer) Scene {
	return &scene{
		world:   donburi.NewWorld(),
		indexer: sceneIndexer,
	}
}

func (s *scene) World() donburi.World {
	return s.world
}

func (s *scene) Render(screen *ebiten.Image, region geom.AABB, view ebiten.GeoM) error {
	// TODO - Query quadtree for renderables.
	return nil
}
