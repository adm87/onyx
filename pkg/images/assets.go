package images

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageAdapter struct {
	logger engine.Logger
	cache  map[engine.FilePath]*ebiten.Image
}

func GetImage(assets engine.Assets, path engine.FilePath) (*ebiten.Image, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	img, exists := adapter.(*ImageAdapter).cache[path]
	return img, exists
}

func NewAdapter(logger engine.Logger) *ImageAdapter {
	return &ImageAdapter{
		logger: logger,
		cache:  make(map[engine.FilePath]*ebiten.Image),
	}
}

func (a *ImageAdapter) GetImage(path engine.FilePath) (*ebiten.Image, bool) {
	img, exists := a.cache[path]
	return img, exists
}

func (a *ImageAdapter) HasImage(path engine.FilePath) bool {
	_, exists := a.cache[path]
	return exists
}

func (a *ImageAdapter) UnloadImage(path engine.FilePath) bool {
	deleted := false

	if img, exists := a.cache[path]; exists {
		img.Deallocate()

		delete(a.cache, path)
		deleted = true
	}

	return deleted
}

func (a *ImageAdapter) ImportAsset(fileSystem fs.FS, path engine.FilePath, raw []byte) error {
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

func (a *ImageAdapter) DeleteAsset(path engine.FilePath) bool {
	return a.UnloadImage(path)
}

func (a *ImageAdapter) SupportedExtensions() []engine.FileExt {
	return []engine.FileExt{".png", ".jpg", ".jpeg"}
}
