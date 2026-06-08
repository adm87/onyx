package images

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageAssetAdapter struct {
}

func NewAdapter() *ImageAssetAdapter {
	return &ImageAssetAdapter{}
}

func (a *ImageAssetAdapter) GetImage(path file.FilePath) (*ebiten.Image, bool) {
	return nil, false
}

func (a *ImageAssetAdapter) HasImage(path file.FilePath) bool {
	return false
}

func (a *ImageAssetAdapter) UnloadImage(path file.FilePath) bool {
	return true
}

func (a *ImageAssetAdapter) ImportAsset(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	return nil
}

func (a *ImageAssetAdapter) DeleteAsset(path file.FilePath) bool {
	return true
}

func (a *ImageAssetAdapter) SupportedExtensions() []file.FileExt {
	return []file.FileExt{".png", ".jpg", ".jpeg"}
}
