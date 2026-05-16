package images

import (
	"bytes"
	"errors"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const adapterID = "images"

type Cache interface {
	Get(path file.Path) (*ebiten.Image, bool)
}

func GetCache(assets engine.Assets) (Cache, bool) {
	adptr, exists := assets.GetAdapter(adapterID)
	if !exists {
		return nil, false
	}

	imageAdapter, ok := adptr.(*adapter)
	if !ok {
		return nil, false
	}

	return imageAdapter, true
}

type adapter struct {
	logger engine.Logger
	images map[file.Path]*ebiten.Image
}

func newAdapter(logger engine.Logger) *adapter {
	return &adapter{
		logger: logger,
		images: make(map[file.Path]*ebiten.Image),
	}
}

func (a *adapter) Import(path file.Path, data []byte) error {
	if len(data) == 0 {
		a.logger.Warn("Received empty data for image at path: %s", path)
		return nil
	}

	if _, exists := a.images[path]; exists {
		return errors.New("image already exists at path: " + string(path))
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		a.logger.Error("Failed to decode image at path: %s, error: %v", path, err)
		return err
	}

	a.images[path] = img
	return nil
}

func (a *adapter) Delete(path file.Path) {
	if img, exists := a.images[path]; exists {
		img.Deallocate()
	}
	delete(a.images, path)
}

func (a *adapter) SupportedTypes() []file.Ext {
	return []file.Ext{".png", ".jpg", ".jpeg"}
}

func (a *adapter) Get(path file.Path) (*ebiten.Image, bool) {
	img, exists := a.images[path]
	return img, exists
}
