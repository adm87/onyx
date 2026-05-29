package images

import (
	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

const AdapterID engine.AdapterID = "image_adapter"

func RegisterPackage(assets engine.Assets, renderer engine.Renderer) error {
	assetAdapter := NewAdapter()
	assets.AddAssetAdapter(
		AdapterID,
		assetAdapter,
	)
	renderer.AddRenderingAdapter(
		AdapterID,
		NewImageRenderingAdapter(
			assetAdapter,
		),
	)
	return nil
}

func GetAssetAdapter(assets engine.Assets) (*ImageAssetAdapter, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	imageAdapter, ok := adapter.(*ImageAssetAdapter)
	return imageAdapter, ok
}

func GetRenderingAdapter(renderer engine.Renderer) (*ImageRenderingAdapter, bool) {
	adapter, found := renderer.GetRenderingAdapter(AdapterID)
	if !found {
		return nil, false
	}

	imageRenderer, ok := adapter.(*ImageRenderingAdapter)
	return imageRenderer, ok
}

func GetImage(assets engine.Assets, path engine.FilePath) (*ebiten.Image, bool) {
	adapter, found := GetAssetAdapter(assets)
	if !found {
		return nil, false
	}

	img, exists := adapter.cache[path]
	return img, exists
}
