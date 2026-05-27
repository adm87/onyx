package components

import (
	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/yohamta/donburi"
)

var Image = donburi.NewComponentType[engine.FilePath]()

func GetImageRef(entry *donburi.Entry) engine.FilePath {
	if !entry.HasComponent(Image) {
		return ""
	}
	return *Image.Get(entry)
}

func SetImageRef(entry *donburi.Entry, ref engine.FilePath) {
	if !entry.HasComponent(Image) {
		entry.AddComponent(Image)
	}
	donburi.SetValue(entry, Image, ref)
}
