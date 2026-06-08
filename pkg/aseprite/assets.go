package aseprite

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
)

const AdapterID engine.AdapterID = "aseprite_adapter"

type AsepriteAssetAdapter struct {
}

func NewAsepriteAssetAdapter() *AsepriteAssetAdapter {
	return &AsepriteAssetAdapter{}
}

func (a *AsepriteAssetAdapter) ImportAsset(filesystem fs.FS, path file.FilePath, raw []byte) error {
	return nil // This asset adapter doesn't handle importing raw data.
}

func (a *AsepriteAssetAdapter) DeleteAsset(path file.FilePath) bool {
	return true // This asset adapter doesn't handle deleting raw data.
}

func (a *AsepriteAssetAdapter) SupportedExtensions() []file.FileExt {
	return nil // This asset adapter doesn't handle importing raw data.
}

func (a *AsepriteAssetAdapter) ImportAnimations(jsonPaths ...file.FilePath) error {
	return nil
}

func (a *AsepriteAssetAdapter) DeleteAnimations(jsonPaths ...file.FilePath) {
}
