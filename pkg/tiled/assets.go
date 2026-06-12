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
	"github.com/adm87/onyx/pkg/images"
)

var tiledExtensions = []file.FileExt{".tmx", ".tsx"}

type assetAdapter struct {
	assets engine.Assets

	tmxStore     *slotmap.SlotMap[*Tmx]
	tsxStore     *slotmap.SlotMap[*Tsx]
	tilemapStore *slotmap.SlotMap[*Tilemap]

	tmxHandles map[file.FilePath]uint64
	tsxHandles map[file.FilePath]uint64

	imageModule *images.ImageModule
}

func newAssetsAdapter(assets engine.Assets, imageModule *images.ImageModule) *assetAdapter {
	return &assetAdapter{
		assets:       assets,
		tmxStore:     slotmap.New[*Tmx](0),
		tsxStore:     slotmap.New[*Tsx](0),
		tilemapStore: slotmap.New[*Tilemap](0),
		tmxHandles:   make(map[file.FilePath]uint64),
		tsxHandles:   make(map[file.FilePath]uint64),
		imageModule:  imageModule,
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
	return true
}

func (a *assetAdapter) SupportedExtensions() []file.FileExt {
	return tiledExtensions
}

func (a *assetAdapter) importTmx(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	if _, exists := a.tmxHandles[path]; exists {
		return nil
	}

	var tmx Tmx

	err := xml.Unmarshal(raw, &tmx)
	assert.Fatal(err)

	tmx.Handle = a.tmxStore.Insert(&tmx)
	a.tmxHandles[path] = tmx.Handle

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

			tileset.Handle = a.tsxHandles[tsxPath]
		}
	}

	return nil
}

func (a *assetAdapter) importTsx(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	if _, exists := a.tsxHandles[path]; exists {
		return nil
	}

	var tsx Tsx

	err := xml.Unmarshal(raw, &tsx)
	assert.Fatal(err)

	tsx.Handle = a.tsxStore.Insert(&tsx)
	a.tsxHandles[path] = tsx.Handle

	dir := filepath.Dir(path.String())

	if tsx.Image.Source == "" {
		return nil
	}

	imgPath := file.ResolvedPath(dir, tsx.Image.Source)
	tsx.Image.Source = imgPath.String()

	err = a.assets.Load(fileSystem, imgPath)
	assert.Fatal(err)

	imgHandle, exists := a.imageModule.GetAssetHandle(imgPath)
	assert.True(exists, fmt.Sprintf("Failed to load image asset for TSX at path %s", imgPath))

	tsx.Image.Handle = imgHandle
	a.imageModule.ExtractUniformFrames(imgHandle, tsx.TileWidth, tsx.TileHeight)

	return nil
}
