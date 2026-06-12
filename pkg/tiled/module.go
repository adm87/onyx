package tiled

import (
	"fmt"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/shapes"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/yohamta/donburi"
)

type TiledModule struct {
	assetsAdapter    *assetAdapter
	renderingAdapter *renderingAdapter

	assetAdapterHandle     uint64
	renderingAdapterHandle uint64
}

func NewModule(
	assets engine.Assets,
	renderer engine.Renderer,
	screen engine.Screen,
	imageModule *images.ImageModule) *TiledModule {

	assetsAdapter := newAssetsAdapter(assets, imageModule)
	renderingAdapter := newRenderingAdapter(screen, imageModule, assetsAdapter)
	return &TiledModule{
		assetsAdapter:          assetsAdapter,
		renderingAdapter:       renderingAdapter,
		assetAdapterHandle:     assets.AddAssetAdapter(assetsAdapter),
		renderingAdapterHandle: renderer.AddRenderingAdapter(renderingAdapter),
	}
}

func (m *TiledModule) BuildTilemap(handle uint64) (*Tilemap, uint64, error) {
	tmx, ok := m.assetsAdapter.tmxStore.Get(handle)
	assert.True(ok, fmt.Sprintf("TMX asset with handle %d not found", handle))

	tilemap, err := buildTilemap(tmx)
	assert.Nil(err, fmt.Sprintf("Failed to build tilemap from TMX asset with handle %d: %v", handle, err))

	return tilemap, m.assetsAdapter.tilemapStore.Insert(tilemap), nil
}

func (m *TiledModule) ReleaseTilemap(handle uint64) {
	m.assetsAdapter.tilemapStore.Delete(handle)
	m.renderingAdapter.releaseBuffer(handle)
}

func (m *TiledModule) GetTilemapSize(handle uint64) (int, int, bool) {
	tilemap, ok := m.assetsAdapter.tilemapStore.Get(handle)
	if !ok {
		return 0, 0, false
	}
	width, height := tilemap.bounds.Width(), tilemap.bounds.Height()
	return int(width), int(height), true
}

func (m *TiledModule) GetTmxHandle(path file.FilePath) (uint64, bool) {
	handle, exists := m.assetsAdapter.tmxHandles[path]
	return handle, exists
}

func (m *TiledModule) GetTsxHandle(path file.FilePath) (uint64, bool) {
	handle, exists := m.assetsAdapter.tsxHandles[path]
	return handle, exists
}

func (m *TiledModule) GetTilemap(handle uint64) (*Tilemap, bool) {
	tilemap, ok := m.assetsAdapter.tilemapStore.Get(handle)
	return tilemap, ok
}

func (m *TiledModule) GetTmx(handle uint64) (*Tmx, bool) {
	tmx, ok := m.assetsAdapter.tmxStore.Get(handle)
	return tmx, ok
}

func (m *TiledModule) GetTsx(handle uint64) (*Tsx, bool) {
	tsx, ok := m.assetsAdapter.tsxStore.Get(handle)
	return tsx, ok
}

func (m *TiledModule) CreateTilemapEntity(ecs donburi.World, opts ...TilemapOption) *donburi.Entry {
	entry := ecs.Entry(ecs.Create(TilemapHandle))

	options := defaultTilemapOptions()
	for _, opt := range opts {
		opt(options)
	}

	SetTilemapHandle(entry, options.Handle)

	width, height, ok := m.GetTilemapSize(options.Handle)
	assert.True(ok, "failed to get tilemap size for the provided handle")

	shapes.AddAABB(entry, shapes.WithBounds(
		geom.Vec2{X: 0, Y: 0},
		geom.Vec2{X: float64(width), Y: float64(height)},
	))

	rendering.AddRenderer(entry,
		rendering.WithRendererID(
			m.renderingAdapterHandle,
		),
	)

	transform.AddTransform(entry)

	return entry
}
