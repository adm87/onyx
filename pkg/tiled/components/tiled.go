package components

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	Tilemap      = donburi.NewComponentType[engine.FilePath]()
	TilemapQuery = donburi.NewQuery(
		filter.Contains(Tilemap),
	)
)

func GetTilemapRef(entry *donburi.Entry) engine.FilePath {
	if !entry.HasComponent(Tilemap) {
		return ""
	}
	return *Tilemap.Get(entry)
}

func SetTilemapRef(entry *donburi.Entry, ref engine.FilePath) {
	if !entry.HasComponent(Tilemap) {
		entry.AddComponent(Tilemap)
	}
	donburi.SetValue(entry, Tilemap, ref)
}
