package images

import (
	"github.com/adm87/onyx/pkg/assert"
	"github.com/adm87/onyx/pkg/engine"
)

const (
	AdapterID engine.AdapterID = "image_adapter"
)

var (
	imageRendererID uint64
)

func RegisterPackage(assets engine.Assets, renderer engine.Renderer) {
	imageAssetAdapter := newAssetAdapter(assets)
	assets.AddAssetAdapter(AdapterID, imageAssetAdapter)

	imageRendererAdapter := newRendererAdapter(imageAssetAdapter)
	imageRendererID = renderer.AddRenderingAdapter(imageRendererAdapter)
}

func GetAssetAdapter(assets engine.Assets) ImageAssetAdapter {
	adapter, ok := assets.GetAdapter(AdapterID)
	assert.True(ok, "image asset adapter not found")

	imageAdapter, ok := adapter.(*assetAdapter)
	assert.True(ok, "image asset adapter type mismatch")

	return imageAdapter
}
