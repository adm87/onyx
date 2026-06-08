package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const (
	AdapterID         engine.AdapterID       = "tiled_adapter"
	TiledRendererType rendering.RendererType = "tiled_renderer"
)

var (
	Tiled      = donburi.NewTag()
	TiledQuery = donburi.NewQuery(
		filter.Contains(Tiled),
	)
)

func RegisterPackage() error {
	return nil
}
