package images

import (
	"bytes"

	"github.com/adm87/onyx/pkg/encoding"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var adapterID = engine.AssetAdapterID(encoding.TypeID[EbitenImageAdapter]())

type EbitenImageAdapter struct {
	logger *engine.Logger
	cache  *cache
}

func NewEbitenImageAdapter(logger *engine.Logger) *EbitenImageAdapter {
	return &EbitenImageAdapter{
		logger: logger,
		cache:  newCache(),
	}
}

func (a *EbitenImageAdapter) SupportedTypes() []engine.FileType {
	return []engine.FileType{"png", "jpeg", "jpg"}
}

func (a *EbitenImageAdapter) ID() engine.AssetAdapterID {
	return adapterID
}

func (a *EbitenImageAdapter) Import(path engine.FilePath, data []byte) error {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		a.logger.Error("failed to import image asset", "path", path, "error", err.Error())
		return err
	}
	a.cache.Set(path, img)
	return nil
}

func (a *EbitenImageAdapter) Delete(path engine.FilePath) error {
	a.cache.Delete(path)
	return nil
}
