package images

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
)

const (
	AdapterID         engine.AdapterID       = "image_adapter"
	ImageRendererType rendering.RendererType = "image_renderer"
)

func RegisterPackage() error {
	return nil
}
