package debug

import (
	"github.com/adm87/onyx/pkg/ecs"
	"github.com/adm87/onyx/pkg/ecs/camera"
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

func DrawTransformBounds(
	donburiECS *ecs.DonburiECS,
	cameraEntry *donburi.Entry,
	target *ebiten.Image,
	screen geom.AABB) {
	viewport := camera.GetViewport(cameraEntry, screen)

	count := 0
	path := vector.Path{}
	donburiECS.QueryAll(viewport, func(entity donburi.Entity) {
		entry := donburiECS.World().Entry(entity)

		bounds := transform.GetWorldBounds(entry)
		if !bounds.Intersects(viewport) {
			return
		}

		min := camera.ToScreen(cameraEntry, bounds.Min, screen)
		max := camera.ToScreen(cameraEntry, bounds.Max, screen)

		path.MoveTo(float32(min.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(max.Y))
		path.LineTo(float32(min.X), float32(max.Y))
		path.Close()

		count++
	})

	println("DrawTransformBounds: ", count)

	vector.StrokePath(target, &path, &vector.StrokeOptions{
		Width: 2,
	}, nil)
}
