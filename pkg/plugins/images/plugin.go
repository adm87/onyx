package images

import (
	"image"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ImagesPlugin struct {
	assetAdapter     *assetAdapter
	renderingAdapter *renderingAdapter

	assetAdapterHandle     uint64
	renderingAdapterHandle uint64
}

func NewImagesPlugin(assets engine.Assets, renderer engine.Renderer) *ImagesPlugin {
	assetAdapter := newAssetAdapter(assets)
	renderingAdapter := newRenderingAdapter(assetAdapter)
	return &ImagesPlugin{
		assetAdapter:           assetAdapter,
		renderingAdapter:       renderingAdapter,
		assetAdapterHandle:     assets.AddAssetAdapter(assetAdapter),
		renderingAdapterHandle: renderer.AddRenderingAdapter(renderingAdapter),
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

	width, height, ok := m.GetFrameSize(options.Handle, options.Frame)
	assert.True(ok, "failed to get image size for the provided handle")

	transform.AddTransform(entry, transform.WithBounds(geom.AABB{
		Min: geom.Vec2{X: 0, Y: 0},
		Max: geom.Vec2{X: float64(width), Y: float64(height)},
	}))

	rendering.AddRenderer(entry,
		rendering.WithRendererID(m.renderingAdapterHandle),
	)

	return entry
}
