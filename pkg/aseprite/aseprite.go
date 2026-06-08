package aseprite

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/asset"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/shapes"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/yohamta/donburi"
)

func RegisterPackage(imageAssetAdapter *images.ImageAssetAdapter, assets engine.Assets, logger engine.Logger, renderer engine.Renderer) error {
	assetAdapter := NewAsepriteAssetAdapter(
		assets,
		logger,
		imageAssetAdapter,
	)
	assets.AddAssetAdapter(
		AdapterID,
		assetAdapter,
	)
	renderer.AddRenderingAdapter(
		AdapterID,
		NewAsepriteRendererAdapter(
			imageAssetAdapter,
			assetAdapter,
		),
	)
	return nil
}

func GetAssetAdapter(assets engine.Assets) (*AsepriteAssetAdapter, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	asepriteAdapter, ok := adapter.(*AsepriteAssetAdapter)
	return asepriteAdapter, ok
}

func CreateSprite(ecs donburi.World, ref file.FilePath, animationName string, bounds geom.AABB) *donburi.Entry {
	entry := asset.NewAssetReference(ecs, ref)
	AddAnimation(entry, animationName)

	rendering.AddRenderer(entry,
		rendering.WithType(AsepriteRendererType),
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
