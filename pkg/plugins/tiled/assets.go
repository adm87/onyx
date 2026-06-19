package tiled

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
	"github.com/adm87/onyx/pkg/plugins/images"
)

var tiledExtensions = []file.FileExt{".tmx", ".tsx"}

type assetAdapter struct {
	assets engine.Assets

	tmxStore file.FileStore[*Tmx]
	tsxStore file.FileStore[*Tsx]

	tilemapStore *slotmap.SlotMap[*Tilemap]
	images       *images.ImagesPlugin
}

func newAssetsAdapter(assets engine.Assets, imagesPlugin *images.ImagesPlugin) *assetAdapter {
	return &assetAdapter{
		assets:       assets,
		tmxStore:     file.NewFileStore[*Tmx](0),
		tsxStore:     file.NewFileStore[*Tsx](0),
		tilemapStore: slotmap.New[*Tilemap](0),
		images:       imagesPlugin,
	}
}

func (a *assetAdapter) ImportAsset(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	switch path.Ext() {
	case ".tmx":
		return a.importTmx(fileSystem, path, raw)
	case ".tsx":
		return a.importTsx(fileSystem, path, raw)
	default:
		return nil
	}
}

func (a *assetAdapter) DeleteAsset(path file.FilePath) bool {
	// TODO - implement this
	return true
}

func (a *assetAdapter) SupportedExtensions() []file.FileExt {
	return tiledExtensions
}

func (a *assetAdapter) importTmx(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	if _, exists := a.tmxStore.GetHandle(path); exists {
		return nil
	}

	var tmx Tmx

	err := xml.Unmarshal(raw, &tmx)
	assert.Fatal(err)

	tmx.Handle = a.tmxStore.Insert(path, &tmx)

	if len(tmx.Tilesets) > 0 {
		dir := filepath.Dir(path.String())

		for i := range tmx.Tilesets {
			tileset := &tmx.Tilesets[i]

			if tileset.Source == "" {
				continue
			}

			tsxPath := file.ResolvedPath(dir, tileset.Source)
			tileset.Source = tsxPath.String()

			err := a.assets.Load(fileSystem, tsxPath)
			assert.Fatal(err)

			handle, exists := a.tsxStore.GetHandle(tsxPath)
			assert.True(exists, fmt.Sprintf("Failed to load TSX asset at path %s", tsxPath))

			tileset.Handle = handle
		}
	}

	return nil
}

func (a *assetAdapter) importTsx(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	if _, exists := a.tsxStore.GetHandle(path); exists {
		return nil
	}

	var tsx Tsx

	err := xml.Unmarshal(raw, &tsx)
	assert.Fatal(err)

	tsx.Handle = a.tsxStore.Insert(path, &tsx)

	dir := filepath.Dir(path.String())

	if tsx.Image.Source == "" {
		return nil
	}

	imgPath := file.ResolvedPath(dir, tsx.Image.Source)
	tsx.Image.Source = imgPath.String()

	err = a.assets.Load(fileSystem, imgPath)
	assert.Fatal(err)

	imgHandle, exists := a.images.GetAssetHandle(imgPath)
	assert.True(exists, fmt.Sprintf("Failed to load image asset for TSX at path %s", imgPath))

	tsx.Image.Handle = imgHandle
	a.images.ExtractUniformFrames(imgHandle, tsx.TileWidth, tsx.TileHeight)

	return nil
}
