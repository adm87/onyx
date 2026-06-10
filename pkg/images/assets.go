package images

import (
	"bytes"
	"io/fs"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageAssetAdapter interface {
	GetImageSize(handle uint64) (int, int, bool)
	GetImage(handle uint64) (*ebiten.Image, bool)
	GetHandle(path file.FilePath) (uint64, bool)
	LoadImage(filesystem fs.FS, path file.FilePath) (uint64, error)
	DeleteImage(path file.FilePath) bool
}

type assetAdapter struct {
	store   *slotmap.SlotMap[*ebiten.Image]
	assets  engine.Assets
	handles map[file.FilePath]uint64
}

func newAssetAdapter(assets engine.Assets) *assetAdapter {
	return &assetAdapter{
		assets:  assets,
		store:   slotmap.New[*ebiten.Image](256),
		handles: make(map[file.FilePath]uint64),
	}
}

func (a *assetAdapter) ImportAsset(_ fs.FS, path file.FilePath, raw []byte) error {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(raw))
	assert.Fatal(err)

	handle := a.store.Insert(img)
	a.handles[path] = handle

	return nil
}

func (a *assetAdapter) DeleteAsset(path file.FilePath) bool {
	handle, exists := a.handles[path]
	if !exists {
		return false
	}

	img, ok := a.store.Delete(handle)
	if !ok {
		return false
	}

	img.Deallocate()
	delete(a.handles, path)

	return true
}

func (a *assetAdapter) SupportedExtensions() []file.FileExt {
	return []file.FileExt{".png", ".jpg", ".jpeg"}
}
