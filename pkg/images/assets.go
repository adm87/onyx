package images

import (
	"bytes"
	"image"
	"io/fs"
	"math"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var imageExtensions = []file.FileExt{".png", ".jpg", ".jpeg"}

type assetAdapter struct {
	store   *slotmap.SlotMap[*ebiten.Image]
	assets  engine.Assets
	handles map[file.FilePath]uint64
	frames  map[uint64][]*ebiten.Image
}

func newAssetAdapter(assets engine.Assets) *assetAdapter {
	return &assetAdapter{
		assets:  assets,
		store:   slotmap.New[*ebiten.Image](256),
		handles: make(map[file.FilePath]uint64),
		frames:  make(map[uint64][]*ebiten.Image),
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

	if frames, exists := a.frames[handle]; exists {
		clear(frames)
		delete(a.frames, handle)
	}

	return true
}

func (a *assetAdapter) SupportedExtensions() []file.FileExt {
	return imageExtensions
}

func (a *assetAdapter) getFrame(handle uint64, index int) (*ebiten.Image, bool) {
	frames, exists := a.frames[handle]
	if !exists || index < 0 || index >= len(frames) {
		if img, exists := a.store.Get(handle); exists {
			return img, true
		}
		return nil, false
	}
	return frames[index], true
}

func (a *assetAdapter) extractUniformFrames(handle uint64, frameWidth, frameHeight int) {
	img, exists := a.store.Get(handle)
	if !exists {
		return
	}

	frames, exists := a.frames[handle]
	if exists {
		frames = frames[:0]
	}

	columns := int(math.Ceil(float64(img.Bounds().Dx()) / float64(frameWidth)))
	rows := int(math.Ceil(float64(img.Bounds().Dy()) / float64(frameHeight)))

	for y := range rows {
		for x := range columns {
			frame := img.SubImage(image.Rect(
				x*frameWidth,
				y*frameHeight,
				(x+1)*frameWidth,
				(y+1)*frameHeight,
			)).(*ebiten.Image)
			frames = append(frames, frame)
		}
	}

	a.frames[handle] = frames
}
