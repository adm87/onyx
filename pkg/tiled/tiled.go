package tiled

import (
	"fmt"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/asset"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/images"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const AdapterID engine.AdapterID = "tiled_adapter"

var (
	Tiled      = donburi.NewTag()
	TiledQuery = donburi.NewQuery(
		filter.Contains(Tiled),
	)
)

func RegisterPackage(assets engine.Assets, renderer engine.Renderer, camera engine.Camera, screen engine.Screen) error {
	imageAssetAdapter, exists := images.GetAssetAdapter(assets)
	if !exists {
		return fmt.Errorf("images asset adapter not found, tiled package requires images package to be registered first")
	}

	tiledAssetAdapter := NewTiledAssetAdapter(imageAssetAdapter)
	assets.AddAssetAdapter(
		AdapterID,
		tiledAssetAdapter,
	)
	renderer.AddRenderingAdapter(
		AdapterID,
		NewTiledRenderingAdapter(
			tiledAssetAdapter,
			imageAssetAdapter,
			camera,
			screen,
		),
	)
	return nil
}

func GetTmx(assets engine.Assets, path file.FilePath) (*Tmx, bool) {
	adapter, found := GetAssetAdapter(assets)
	if !found {
		return nil, false
	}

	tmx, exists := adapter.tmxCache[path]
	return tmx, exists
}

func GetTsx(assets engine.Assets, path file.FilePath) (*Tsx, bool) {
	adapter, found := GetAssetAdapter(assets)
	if !found {
		return nil, false
	}

	tsx, exists := adapter.tsxCache[path]
	return tsx, exists
}

func GetTilemap(assets engine.Assets, path file.FilePath) (*Tilemap, bool) {
	adapter, found := GetAssetAdapter(assets)
	if !found {
		return nil, false
	}

	tilemap, exists := adapter.tilemaps[path]
	return tilemap, exists
}

func GetAssetAdapter(assets engine.Assets) (*TiledAssetAdapter, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	tiledAdapter, ok := adapter.(*TiledAssetAdapter)
	return tiledAdapter, ok
}

func GetRenderingAdapter(renderer engine.Renderer) (*TiledRenderingAdapter, bool) {
	adapter, found := renderer.GetRenderingAdapter(AdapterID)
	if !found {
		return nil, false
	}

	tiledRenderer, ok := adapter.(*TiledRenderingAdapter)
	return tiledRenderer, ok
}

func CreateTiledEntity(ecs donburi.World, ref file.FilePath) *donburi.Entry {
	entry := asset.NewAssetReference(ecs, ref)
	entry.AddComponent(Tiled)

	transform.AddTransform(entry)
	rendering.AddRenderer(entry)

	return entry
}
