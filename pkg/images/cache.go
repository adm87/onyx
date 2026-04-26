package images

import (
	"sync"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

func Cache(assets *engine.Assets) (engine.AssetCache[*ebiten.Image], error) {
	adapter, exists := assets.AdapterByID(adapterID)
	if !exists {
		return nil, engine.ErrAssetAdapterNotFound
	}
	return adapter.(*EbitenImageAdapter).cache, nil
}

type cache struct {
	images map[engine.FilePath]*ebiten.Image
	mu     sync.RWMutex
}

func newCache() *cache {
	return &cache{
		images: make(map[engine.FilePath]*ebiten.Image),
	}
}

func (c *cache) Get(path engine.FilePath) (*ebiten.Image, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	img, exists := c.images[path]
	return img, exists
}

func (c *cache) Set(path engine.FilePath, img *ebiten.Image) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.images[path] = img
}

func (c *cache) Delete(path engine.FilePath) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if img, exists := c.images[path]; exists {
		img.Deallocate()
		delete(c.images, path)
	}
}
