package images

import (
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	Components = []donburi.IComponentType{
		transform.Position,
		transform.Rotation,
		transform.Scale,
		transform.Matrix,
		rendering.Renderer,
		rendering.Image,
		rendering.Anchor,
		rendering.Color,
	}
)

var (
	Query = donburi.NewQuery(
		filter.Contains(Components...),
	)
)
