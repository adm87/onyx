package tiled

import (
	"fmt"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
)

var (
	AdapterID engine.AdapterID = "tiled_adapter"
)

func RegisterPackage(assets engine.Assets, renderer engine.Renderer, logger engine.Logger) error {
	logger.Debug("Registering tiled package")

	imageAssetAdapter, err := getImageAssetAdapter(assets)
	if err != nil {
		return fmt.Errorf("failed to get image adapter: %w", err)
	}

	assets.AddAssetAdapter(
		AdapterID,
		NewTiledAssetAdapter(
			imageAssetAdapter,
			logger,
		),
	)

	renderer.AddRenderingAdapter(
		AdapterID,
		NewTiledRenderingAdapter(
			logger,
		),
	)

	return nil
}

func getImageAssetAdapter(assets engine.Assets) (*images.ImageAdapter, error) {
	adapter, found := assets.GetAdapter(images.AdapterID)
	if !found {
		return nil, fmt.Errorf("image adapter with ID '%s' not found", images.AdapterID)
	}

	imageAdapter, ok := adapter.(*images.ImageAdapter)
	if !ok {
		return nil, fmt.Errorf("adapter with ID '%s' is not of type '*images.ImageAdapter'", images.AdapterID)
	}

	return imageAdapter, nil
}
