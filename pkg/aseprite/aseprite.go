package aseprite

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
)

func RegisterPackage(imageAssetAdapter *images.ImageAssetAdapter, assets engine.Assets, renderer engine.Renderer) error {
	assetAdapter := NewAsepriteAssetAdapter(imageAssetAdapter)
	assets.AddAssetAdapter(
		AdapterID,
		assetAdapter,
	)
	return nil
}

func GetAssetAdapter(assets engine.Assets) (*AsepriteAssetAdapter, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	asepriteAdapter, ok := adapter.(*AsepriteAssetAdapter)
	return asepriteAdapter, ok
}
