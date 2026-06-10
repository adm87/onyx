package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ImageModule struct {
	assetAdapter    *assetAdapter
	rendererAdapter *renderingAdapter

	assetAdapterHandle     uint64
	renderingAdapterHandle uint64
}

func NewModule(assets engine.Assets, renderer engine.Renderer) *ImageModule {
	assetAdapter := newAssetAdapter(assets)
	renderingAdapter := newRendererAdapter(assetAdapter)
	return &ImageModule{
		assetAdapter:           assetAdapter,
		rendererAdapter:        renderingAdapter,
		assetAdapterHandle:     assets.AddAssetAdapter(assetAdapter),
		renderingAdapterHandle: renderer.AddRenderingAdapter(renderingAdapter),
	}
}

func (m *ImageModule) GetAssetHandle(path file.FilePath) (uint64, bool) {
	handle, exists := m.assetAdapter.handles[path]
	return handle, exists
}

func (m *ImageModule) GetImageSize(handle uint64) (int, int, bool) {
	img, exists := m.assetAdapter.store.Get(handle)
	if !exists {
		return 0, 0, false
	}
	return img.Bounds().Dx(), img.Bounds().Dy(), true
}

func (m *ImageModule) GetImage(handle uint64) (*ebiten.Image, bool) {
	return m.assetAdapter.store.Get(handle)
}

func (m *ImageModule) CreateImage(ecs donburi.World, opts ...ImageOption) *donburi.Entry {
	entry := ecs.Entry(ecs.Create(ImageHandle))

	options := defaultImageOptions()
	for _, opt := range opts {
		opt(options)
	}

	SetImageHandle(entry, options.Handle)

	transform.AddTransform(entry)
	rendering.AddRenderer(entry, rendering.WithRendererID(m.renderingAdapterHandle))

	return entry
}
