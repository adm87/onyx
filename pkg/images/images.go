package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/asset"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/shapes"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const (
	AdapterID         engine.AdapterID       = "image_adapter"
	ImageRendererType rendering.RendererType = "image_renderer"
)

var (
	Image      = donburi.NewTag()
	ImageQuery = donburi.NewQuery(
		filter.Contains(Image),
	)
)

func RegisterPackage(assets engine.Assets, renderer engine.Renderer) error {
	assetAdapter := NewAdapter()
	assets.AddAssetAdapter(
		AdapterID,
		assetAdapter,
	)
	renderer.AddRenderingAdapter(
		AdapterID,
		NewImageRenderingAdapter(
			assetAdapter,
		),
	)
	return nil
}

func GetAssetAdapter(assets engine.Assets) (*ImageAssetAdapter, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	imageAdapter, ok := adapter.(*ImageAssetAdapter)
	return imageAdapter, ok
}

func GetRenderingAdapter(renderer engine.Renderer) (*ImageRenderingAdapter, bool) {
	adapter, found := renderer.GetRenderingAdapter(AdapterID)
	if !found {
		return nil, false
	}

	imageRenderer, ok := adapter.(*ImageRenderingAdapter)
	return imageRenderer, ok
}

func GetImageAssets(assets engine.Assets, path file.FilePath) (*ebiten.Image, bool) {
	adapter, found := GetAssetAdapter(assets)
	if !found {
		return nil, false
	}

	img, exists := adapter.cache[path]
	return img, exists
}

func CreateImageEntity(ecs donburi.World, ref file.FilePath, bounds geom.AABB) *donburi.Entry {
	entry := asset.NewAssetReference(ecs, ref)
	entry.AddComponent(Image)

	rendering.AddRenderer(entry,
		rendering.WithType(ImageRendererType),
	)
	shapes.AddAABB(entry,
		shapes.WithBounds(
			bounds.Min,
			bounds.Max,
		),
	)
	transform.AddTransform(entry)

	return entry
}
