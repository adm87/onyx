package game

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
)

func addAssetAdapters(onyx engine.Game) {
	assets := onyx.Assets()

	assets.AddAdapter(images.AdapterID, images.NewAdapter())
	assets.AddAdapter(tiled.AdapterID, tiled.NewAdapter())
}
