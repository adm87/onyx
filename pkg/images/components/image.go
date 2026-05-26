package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ImageReg struct {
	Ref *ebiten.Image
}

var Image = donburi.NewComponentType[ImageReg]()

func GetImage(entry *donburi.Entry) *ebiten.Image {
	if !entry.HasComponent(Image) {
		return nil
	}
	return Image.Get(entry).Ref
}

func SetImage(entry *donburi.Entry, ref *ebiten.Image) {
	if !entry.HasComponent(Image) {
		entry.AddComponent(Image)
	}
	donburi.SetValue(entry, Image, ImageReg{Ref: ref})
}
