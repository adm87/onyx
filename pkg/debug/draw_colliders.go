package debug

import (
	"image/color"

	"github.com/adm87/onyx/pkg/ecs/camera"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

func DrawColliders(
	world donburi.World,
	collisionWorld *collision.CollisionWorld,
	cameraEntry *donburi.Entry,
	target *ebiten.Image,
	screen geom.AABB) {
	viewport := camera.GetViewport(cameraEntry, screen)

	path := vector.Path{}
	collisionWorld.QueryAll(viewport, func(entity donburi.Entity) {
		entry := world.Entry(entity)

		collider := collision.GetWorldCollider(entry)
		if !collider.Intersects(viewport) {
			return
		}

		min := camera.ToScreen(cameraEntry, collider.Min, screen)
		max := camera.ToScreen(cameraEntry, collider.Max, screen)

		path.MoveTo(float32(min.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(max.Y))
		path.LineTo(float32(min.X), float32(max.Y))
		path.Close()
	})

	opt := &vector.DrawPathOptions{}
	opt.ColorScale.ScaleWithColor(color.RGBA{R: 255, A: 255})

	vector.StrokePath(target, &path, &vector.StrokeOptions{
		Width: 2,
	}, opt)
}
