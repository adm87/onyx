package images

import (
	"image"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ImagesPlugin struct {
	assetAdapter *assetAdapter

	assetAdapterHandle     uint64
	renderingAdapterHandle uint64
}

func NewImagesPlugin(assets engine.Assets, renderer engine.Renderer) *ImagesPlugin {
	assetAdapter := newAssetAdapter(assets)
	return &ImagesPlugin{
		assetAdapter:       assetAdapter,
		assetAdapterHandle: assets.AddAssetAdapter(assetAdapter),
	}
}

func (m *ImagesPlugin) GetAssetHandle(path file.FilePath) (uint64, bool) {
	return m.assetAdapter.store.GetHandle(path)
}

func (m *ImagesPlugin) GetImageSize(handle uint64) (int, int, bool) {
	img, exists := m.assetAdapter.store.Get(handle)
	if !exists {
		return 0, 0, false
	}
	return img.Bounds().Dx(), img.Bounds().Dy(), true
}

func (m *ImagesPlugin) GetFrameSize(handle uint64, frame int) (int, int, bool) {
	img, exists := m.assetAdapter.getFrame(handle, frame)
	if !exists {
		return 0, 0, false
	}
	return img.Bounds().Dx(), img.Bounds().Dy(), true
}

func (m *ImagesPlugin) GetImage(handle uint64) (*ebiten.Image, bool) {
	return m.assetAdapter.store.Get(handle)
}

func (m *ImagesPlugin) GetFrameImage(handle uint64, index int) (*ebiten.Image, bool) {
	return m.assetAdapter.getFrame(handle, index)
}

func (m *ImagesPlugin) ExtractUniformFrames(handle uint64, frameWidth, frameHeight int) {
	m.assetAdapter.extractUniformFrames(handle, frameWidth, frameHeight)
}

func (m *ImagesPlugin) ExtractFrames(handle uint64, rects []image.Rectangle) {
	m.assetAdapter.extractFrames(handle, rects)
}

func (m *ImagesPlugin) CreateImageEntity(ecs donburi.World, opts ...Option) *donburi.Entry {
	options := defaultImageOptions()
	for _, opt := range opts {
		opt(options)
	}

	entry := ecs.Entry(ecs.Create(
		Image,
	))

	img := GetImage(entry)
	img.Anchor = options.Anchor
	img.Filter = options.Filter
	img.Frame = options.Frame
	img.Handle = options.Handle
	img.Color = options.Color

	return entry
}
