package aseprite

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/images"
)

const AdapterID engine.AdapterID = "aseprite_adapter"

type AsepriteAssetAdapter struct {
	imageAssetAdapter *images.ImageAssetAdapter
}

func NewAsepriteAssetAdapter(imageAssetAdapter *images.ImageAssetAdapter) *AsepriteAssetAdapter {
	return &AsepriteAssetAdapter{
		imageAssetAdapter: imageAssetAdapter,
	}
}

func (a *AsepriteAssetAdapter) ImportAsset(filesystem fs.FS, path file.FilePath, raw []byte) error {
	return nil
}

func (a *AsepriteAssetAdapter) DeleteAsset(path file.FilePath) bool {
	return true
}

func (a *AsepriteAssetAdapter) SupportedExtensions() []file.FileExt {
	return nil // This asset adapter doesn't handle importing raw data.
}

func (a *AsepriteAssetAdapter) ImportAnimations(jsonPaths ...file.FilePath) error {
	return nil
}

func (a *AsepriteAssetAdapter) DeleteAnimations(jsonPaths ...file.FilePath) {

}
