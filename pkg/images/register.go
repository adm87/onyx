package images

import "github.com/adm87/onyx/pkg/engine"

func Register(assets engine.Assets, logger engine.Logger) {
	assets.RegisterAdapter("images", newAdapter(logger))
}
