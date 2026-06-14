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

type ImageModule struct {
	assetAdapter     *assetAdapter
	renderingAdapter *renderingAdapter

	assetAdapterHandle     uint64
	renderingAdapterHandle uint64
}

func NewModule(assets engine.Assets, renderer engine.Renderer) *ImageModule {
	assetAdapter := newAssetAdapter(assets)
	renderingAdapter := newRenderingAdapter(assetAdapter)
	return &ImageModule{
		assetAdapter:           assetAdapter,
		renderingAdapter:       renderingAdapter,
		assetAdapterHandle:     assets.AddAssetAdapter(assetAdapter),
		renderingAdapterHandle: renderer.AddRenderingAdapter(renderingAdapter),
	}
}

func (m *ImageModule) GetAssetHandle(path file.FilePath) (uint64, bool) {
	return m.assetAdapter.store.GetHandle(path)
}

func (m *ImageModule) GetImageSize(handle uint64) (int, int, bool) {
	img, exists := m.assetAdapter.store.Get(handle)
	if !exists {
		return 0, 0, false
	}
	return img.Bounds().Dx(), img.Bounds().Dy(), true
}

func (m *ImageModule) GetFrameSize(handle uint64, frame int) (int, int, bool) {
	img, exists := m.assetAdapter.getFrame(handle, frame)
	if !exists {
		return 0, 0, false
	}
	return img.Bounds().Dx(), img.Bounds().Dy(), true
}

func (m *ImageModule) GetImage(handle uint64) (*ebiten.Image, bool) {
	return m.assetAdapter.store.Get(handle)
}

func (m *ImageModule) GetFrameImage(handle uint64, index int) (*ebiten.Image, bool) {
	return m.assetAdapter.getFrame(handle, index)
}

func (m *ImageModule) ExtractUniformFrames(handle uint64, frameWidth, frameHeight int) {
	m.assetAdapter.extractUniformFrames(handle, frameWidth, frameHeight)
}

func (m *ImageModule) ExtractFrames(handle uint64, rects []image.Rectangle) {
	m.assetAdapter.extractFrames(handle, rects)
}

func (m *ImageModule) CreateImageEntity(ecs donburi.World, opts ...Option) *donburi.Entry {
	options := defaultImageOptions()
	for _, opt := range opts {
		opt(options)
	}

	entry := ecs.Entry(ecs.Create(
		Image,
		Frame,
	))

	SetImageHandle(entry, options.Handle)
	SetFrame(entry, options.Frame)

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
