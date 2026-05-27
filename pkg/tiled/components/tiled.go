package components

import (
	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/yohamta/donburi"
)

var Tiled = donburi.NewComponentType[engine.FilePath]()

func GetTiledRef(entry *donburi.Entry) engine.FilePath {
	if !entry.HasComponent(Tiled) {
		return ""
	}
	return *Tiled.Get(entry)
}

func SetTiledRef(entry *donburi.Entry, ref engine.FilePath) {
	if !entry.HasComponent(Tiled) {
		entry.AddComponent(Tiled)
	}
	donburi.SetValue(entry, Tiled, ref)
}
