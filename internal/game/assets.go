package game

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
)

func addAssetAdapters(onyx engine.Game) {
	assets := onyx.Assets()
	logger := onyx.Logger()

	imageAdapter := images.NewAdapter()
	assets.AddAdapter(
		images.AdapterID,
		imageAdapter,
	)

	assets.AddAdapter(
		tiled.AdapterID,
		tiled.NewAdapter(imageAdapter, logger),
	)
}
