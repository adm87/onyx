package tiled

import (
	"fmt"

	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/adm87/onyx-game/pkg/images"
	"github.com/adm87/onyx-game/pkg/tiled/components"
	"github.com/adm87/onyx-game/pkg/tiled/data"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const AdapterID engine.AdapterID = "tiled_adapter"

var TiledQuery = donburi.NewQuery(
	filter.Contains(components.Tiled),
)

func RegisterPackage(assets engine.Assets, renderer engine.Renderer, screen engine.Screen) error {
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
			screen,
		),
	)

	return nil
}

func GetTmx(assets engine.Assets, path engine.FilePath) (*data.Tmx, bool) {
	adapter, found := GetAssetAdapter(assets)
	if !found {
		return nil, false
	}

	tmx, exists := adapter.tmxCache[path]
	return tmx, exists
}

func GetTsx(assets engine.Assets, path engine.FilePath) (*data.Tsx, bool) {
	adapter, found := GetAssetAdapter(assets)
	if !found {
		return nil, false
	}

	tsx, exists := adapter.tsxCache[path]
	return tsx, exists
}

func GetTilemap(assets engine.Assets, path engine.FilePath) (*Tilemap, bool) {
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
