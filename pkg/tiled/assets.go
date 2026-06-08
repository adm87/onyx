package tiled

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine/file"
)

type TiledAssetAdapter struct {
}

func NewTiledAssetAdapter() *TiledAssetAdapter {
	return &TiledAssetAdapter{}
}

func (a *TiledAssetAdapter) ImportAsset(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	return nil
}

func (a *TiledAssetAdapter) DeleteAsset(path file.FilePath) bool {
	return true
}

func (a *TiledAssetAdapter) SupportedExtensions() []file.FileExt {
	return []file.FileExt{".tmx", ".tsx"}
}
