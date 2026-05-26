package tiled

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine"
)

var AdapterID engine.AssetAdapterID = "tiled"

type assetAdapter struct {
}

func NewAdapter() *assetAdapter {
	return &assetAdapter{}
}

func (a *assetAdapter) ImportAsset(fileSystem fs.FS, path engine.FilePath, raw []byte) error {
	return nil
}

func (a *assetAdapter) DeleteAsset(path engine.FilePath) bool {
	return false
}

func (a *assetAdapter) SupportedExtensions() []engine.FileExt {
	return []engine.FileExt{".tmx"}
}
