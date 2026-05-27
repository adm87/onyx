package images

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageAssetAdapter struct {
	cache map[engine.FilePath]*ebiten.Image
}

func NewAdapter() *ImageAssetAdapter {
	return &ImageAssetAdapter{
		cache: make(map[engine.FilePath]*ebiten.Image),
	}
}

func (a *ImageAssetAdapter) GetImage(path engine.FilePath) (*ebiten.Image, bool) {
	img, exists := a.cache[path]
	return img, exists
}

func (a *ImageAssetAdapter) HasImage(path engine.FilePath) bool {
	_, exists := a.cache[path]
	return exists
}

func (a *ImageAssetAdapter) UnloadImage(path engine.FilePath) bool {
	deleted := false

	if img, exists := a.cache[path]; exists {
		img.Deallocate()

		delete(a.cache, path)
		deleted = true
	}

	return deleted
}

func (a *ImageAssetAdapter) ImportAsset(fileSystem fs.FS, path engine.FilePath, raw []byte) error {
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

func (a *ImageAssetAdapter) DeleteAsset(path engine.FilePath) bool {
	return a.UnloadImage(path)
}

func (a *ImageAssetAdapter) SupportedExtensions() []engine.FileExt {
	return []engine.FileExt{".png", ".jpg", ".jpeg"}
}
