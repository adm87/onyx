package images

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const AdapterID engine.AssetAdapterID = "images"

type assetAdapter struct {
	cache map[engine.FilePath]*ebiten.Image
}

func GetImage(assets engine.Assets, path engine.FilePath) (*ebiten.Image, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	img, exists := adapter.(*assetAdapter).cache[path]
	return img, exists
}

func NewAdapter() *assetAdapter {
	return &assetAdapter{
		cache: make(map[engine.FilePath]*ebiten.Image),
	}
}

func (a *assetAdapter) ImportAsset(fileSystem fs.FS, path engine.FilePath, raw []byte) error {
	if _, exists := a.cache[path]; exists {
		return fmt.Errorf("asset with path '%s' already exists", path)
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("failed to decode image asset '%s': %w", path, err)
	}

	a.cache[path] = img
	return nil
}

func (a *assetAdapter) DeleteAsset(path engine.FilePath) bool {
	deleted := false

	if img, exists := a.cache[path]; exists {
		img.Deallocate()
		deleted = true
	}
	delete(a.cache, path)

	return deleted
}

func (a *assetAdapter) SupportedExtensions() []engine.FileExt {
	return []engine.FileExt{".png", ".jpg", ".jpeg"}
}
