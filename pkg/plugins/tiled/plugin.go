package tiled

import (
	"fmt"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/yohamta/donburi"
)

type TiledPlugin struct {
	assetsAdapter *assetAdapter

	assetAdapterHandle     uint64
	renderingAdapterHandle uint64
}

func NewTiledPlugin(
	assets engine.Assets,
	renderer engine.Renderer,
	screen engine.Screen,
	imagesPlugin *images.ImagesPlugin) *TiledPlugin {

	assetsAdapter := newAssetsAdapter(assets, imagesPlugin)
	return &TiledPlugin{
		assetsAdapter:      assetsAdapter,
		assetAdapterHandle: assets.AddAssetAdapter(assetsAdapter),
	}
}

func (m *TiledPlugin) BuildTilemap(handle uint64) (*Tilemap, uint64) {
	tmx, ok := m.assetsAdapter.tmxStore.Get(handle)
	assert.True(ok, fmt.Sprintf("TMX asset with handle %d not found", handle))

	tilemap, err := buildTilemap(tmx)
	assert.Nil(err, fmt.Sprintf("Failed to build tilemap from TMX asset with handle %d: %v", handle, err))

	return tilemap, m.assetsAdapter.tilemapStore.Insert(tilemap)
}

func (m *TiledPlugin) ReleaseTilemap(handle uint64) {
	m.assetsAdapter.tilemapStore.Delete(handle)
}

func (m *TiledPlugin) GetTilemapSize(handle uint64) (int, int, bool) {
	tilemap, ok := m.assetsAdapter.tilemapStore.Get(handle)
	if !ok {
		return 0, 0, false
	}
	width, height := tilemap.bounds.Width(), tilemap.bounds.Height()
	return int(width), int(height), true
}

func (m *TiledPlugin) GetTmxHandle(path file.FilePath) (uint64, bool) {
	return m.assetsAdapter.tmxStore.GetHandle(path)
}

func (m *TiledPlugin) GetTsxHandle(path file.FilePath) (uint64, bool) {
	return m.assetsAdapter.tsxStore.GetHandle(path)
}

func (m *TiledPlugin) GetTilemap(handle uint64) (*Tilemap, bool) {
	tilemap, ok := m.assetsAdapter.tilemapStore.Get(handle)
	return tilemap, ok
}

func (m *TiledPlugin) GetTmx(handle uint64) (*Tmx, bool) {
	tmx, ok := m.assetsAdapter.tmxStore.Get(handle)
	return tmx, ok
}

func (m *TiledPlugin) GetTsx(handle uint64) (*Tsx, bool) {
	tsx, ok := m.assetsAdapter.tsxStore.Get(handle)
	return tsx, ok
}

func (m *TiledPlugin) CreateTilemapEntity(ecs donburi.World, opts ...TilemapOption) *donburi.Entry {
	entry := ecs.Entry(ecs.Create(TilemapHandle))

	options := defaultTilemapOptions()
	for _, opt := range opts {
		opt(options)
	}

	SetTilemapHandle(entry, options.Handle)

	return entry
}
