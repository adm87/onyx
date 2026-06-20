package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/images"
)

type TiledPlugin struct {
	logger   engine.Logger
	assets   *TiledAssets
	tilemaps map[uint64]*Tilemap
}

func NewTiledPlugin(logger engine.Logger, images *images.ImageAssets) *TiledPlugin {
	return &TiledPlugin{
		logger:   logger,
		assets:   NewTiledAssets(images),
		tilemaps: make(map[uint64]*Tilemap),
	}
}

func (t *TiledPlugin) Assets() *TiledAssets {
	return t.assets
}

func (t *TiledPlugin) BuildTilemap(handle uint64) *Tilemap {
	tmx, ok := t.assets.tmxStore.Get(handle)
	if !ok {
		return nil
	}

	tilemap, err := buildTilemap(tmx)
	if err != nil {
		return nil
	}

	t.tilemaps[handle] = tilemap
	return tilemap
}

func (t *TiledPlugin) DeleteTilemap(handle uint64) {
	delete(t.tilemaps, handle)
}

func (t *TiledPlugin) GetTilemap(handle uint64) (*Tilemap, bool) {
	tilemap, ok := t.tilemaps[handle]
	return tilemap, ok
}
