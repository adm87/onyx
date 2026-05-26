package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const AdapterID engine.AdapterID = "image_adapter"

var ImageQuery = donburi.NewQuery(
	filter.Contains(components.Image),
)

func RegisterPackage(assets engine.Assets, renderer engine.Renderer, logger engine.Logger) error {
	logger.Debug("Registering images package")

	assets.AddAssetAdapter(
		AdapterID,
		NewAdapter(),
	)

	renderer.AddRenderingAdapter(
		AdapterID,
		NewImageRenderingAdapter(),
	)

	return nil
}
