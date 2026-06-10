package spatial

import "github.com/yohamta/donburi"

var spatialIndexing = donburi.NewComponentType[uint64]()

func AddSpatialIndexing(entry *donburi.Entry, handle uint64) {
	if entry.HasComponent(spatialIndexing) {
		return
	}
	donburi.Add(entry, spatialIndexing, &handle)
}

func GetSpatialIndexing(entry *donburi.Entry) (uint64, bool) {
	if !entry.HasComponent(spatialIndexing) {
		return 0, false
	}
	return *spatialIndexing.Get(entry), true
}
