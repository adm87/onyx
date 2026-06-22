package camera

import (
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var MainCamera = donburi.NewTag("Main Camera")

func GetMainCamera(world donburi.World) (*donburi.Entry, bool) {
	return MainCamera.First(world)
}

func GetViewport(entry *donburi.Entry, screenBounds geom.AABB) geom.AABB {
	viewMatrix := GetViewMatrix(entry, screenBounds)
	viewMatrix.Invert()

	minX, minY := viewMatrix.Apply(screenBounds.Min.X, screenBounds.Min.Y)
	maxX, maxY := viewMatrix.Apply(screenBounds.Max.X, screenBounds.Max.Y)

	return geom.AABB{
		Min: geom.Vec2{X: minX, Y: minY},
		Max: geom.Vec2{X: maxX, Y: maxY},
	}
}

func GetViewMatrix(entry *donburi.Entry, screenBounds geom.AABB) ebiten.GeoM {
	width, height := screenBounds.Width(), screenBounds.Height()

	matrix := transform.GetMatrix(entry)
	matrix.Invert()

	x, y := (width*0.5)+screenBounds.Min.X, (height*0.5)+screenBounds.Min.Y
	matrix.Translate(x, y)

	return matrix
}

func GetZoom(entry *donburi.Entry) float64 {
	x, y := transform.GetScale(entry)
	return (x + y) * 0.5
}

func SetZoom(entry *donburi.Entry, zoom float64) {
	transform.SetScale(entry, zoom, zoom)
}

func ToWorld(entry *donburi.Entry, position geom.Vec2, screenBounds geom.AABB) geom.Vec2 {
	viewMatrix := GetViewMatrix(entry, screenBounds)
	viewMatrix.Invert()

	worldX, worldY := viewMatrix.Apply(position.X, position.Y)
	return geom.Vec2{X: worldX, Y: worldY}
}

func ToScreen(entry *donburi.Entry, position geom.Vec2, screenBounds geom.AABB) geom.Vec2 {
	viewMatrix := GetViewMatrix(entry, screenBounds)

	screenX, screenY := viewMatrix.Apply(position.X, position.Y)
	return geom.Vec2{X: screenX, Y: screenY}
}
