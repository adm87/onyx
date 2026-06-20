package debug

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/ecs"
	"github.com/adm87/onyx/pkg/plugins/ecs/camera"
	"github.com/adm87/onyx/pkg/plugins/ecs/transform"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

func DrawTransformBounds(ecsPlugin *ecs.DonburiECSPlugin, cameraEntry *donburi.Entry, target *ebiten.Image, screen geom.AABB) {
	viewport := camera.GetViewport(cameraEntry, screen)

	path := vector.Path{}
	ecsPlugin.Query(viewport, func(entity donburi.Entity) {
		entry := ecsPlugin.World().Entry(entity)

		bounds := transform.GetWorldBounds(entry)

		min := camera.ToScreen(cameraEntry, bounds.Min, screen)
		max := camera.ToScreen(cameraEntry, bounds.Max, screen)

		path.MoveTo(float32(min.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(max.Y))
		path.LineTo(float32(min.X), float32(max.Y))
		path.Close()
	})

	vector.StrokePath(target, &path, &vector.StrokeOptions{
		Width: 2,
	}, nil)
}
