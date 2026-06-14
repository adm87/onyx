package scene

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
)

var (
	SceneIndexing = donburi.NewComponentType[uint64]()
	SceneBounds   = donburi.NewComponentType[geom.AABB]()
)

func AddSceneIndexing(entry *donburi.Entry, handle uint64) {
	if entry.HasComponent(SceneIndexing) {
		return
	}
	donburi.Add(entry, SceneIndexing, &handle)
}

func GetSceneIndexing(entry *donburi.Entry) (uint64, bool) {
	if !entry.HasComponent(SceneIndexing) {
		return 0, false
	}
	return *SceneIndexing.Get(entry), true
}

func AddSceneBounds(entry *donburi.Entry, bounds *geom.AABB) {
	if entry.HasComponent(SceneBounds) {
		return
	}
	donburi.Add(entry, SceneBounds, bounds)
}

func GetSceneBounds(entry *donburi.Entry) *geom.AABB {
	if !entry.HasComponent(SceneBounds) {
		return &geom.AABB{
			Min: geom.Vec2{X: 0, Y: 0},
			Max: geom.Vec2{X: 0, Y: 0},
		}
	}
	return SceneBounds.Get(entry)
}

func SetSceneBounds(entry *donburi.Entry, bounds *geom.AABB) {
	donburi.Add(entry, SceneBounds, bounds)
}
