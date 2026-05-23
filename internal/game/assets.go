package game

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
)

func addAssetAdapters(onyx engine.Game) {
	assets := onyx.Assets()

	assets.AddAdapter(images.AdapterID, images.NewAdapter())
}
