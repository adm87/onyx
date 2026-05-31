package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	TilemapComponent = donburi.NewComponentType[engine.FilePath]()
	TilemapQuery     = donburi.NewQuery(
		filter.Contains(TilemapComponent),
	)
)

func GetTilemapRef(entry *donburi.Entry) engine.FilePath {
	if !entry.HasComponent(TilemapComponent) {
		return ""
	}
	return *TilemapComponent.Get(entry)
}

func SetTilemapRef(entry *donburi.Entry, ref engine.FilePath) {
	if !entry.HasComponent(TilemapComponent) {
		entry.AddComponent(TilemapComponent)
	}
	donburi.SetValue(entry, TilemapComponent, ref)
}
