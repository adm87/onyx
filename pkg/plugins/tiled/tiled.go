package tiled

import (
	"github.com/adm87/onyx/pkg/plugins/images"
)

type TiledPlugin struct {
	assets *TiledAssets
}

func NewTiledPlugin(images *images.ImageAssets) *TiledPlugin {
	return &TiledPlugin{
		assets: NewTiledAssets(images),
	}
}

func (t *TiledPlugin) Assets() *TiledAssets {
	return t.assets
}
