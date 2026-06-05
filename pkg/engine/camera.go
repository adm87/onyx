package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var CameraTag = donburi.NewTag("Camera")

// Camera defines the interface for interacting with a singleton camera entity within the game world.
// It provides methods to get and set the camera's position and zoom level, as well as to convert between world and screen coordinates.
type Camera interface {
	Position(ecs donburi.World) geom.Vec2
	SetPosition(ecs donburi.World, pos geom.Vec2)
	Zoom(ecs donburi.World) float64
	SetZoom(ecs donburi.World, zoom float64)
	ToWorld(ecs donburi.World, screen Screen, position geom.Vec2) geom.Vec2
	ToScreen(ecs donburi.World, screen Screen, position geom.Vec2) geom.Vec2
}

type camera struct {
	donburi.Entity
}

func newCamera(ecs donburi.World) *camera {
	return &camera{transform.NewTransform(ecs).Entity()}
}

func (c *camera) Position(ecs donburi.World) geom.Vec2 {
	entry := ecs.Entry(c.Entity)
	return transform.GetPosition(entry)
}

func (c *camera) SetPosition(ecs donburi.World, pos geom.Vec2) {
	entry := ecs.Entry(c.Entity)
	transform.SetPosition(entry, pos)
}

func (c *camera) Zoom(ecs donburi.World) float64 {
	entry := ecs.Entry(c.Entity)
	scale := transform.GetScale(entry)
	return scale.X
}

func (c *camera) SetZoom(ecs donburi.World, zoom float64) {
	entry := ecs.Entry(c.Entity)
	transform.SetScale(entry, geom.Vec2{X: zoom, Y: zoom})
}

func (c *camera) view(ecs donburi.World, screen Screen) ebiten.GeoM {
	entry := ecs.Entry(c.Entity)
	safeArea := screen.SafeArea()

	matrix := transform.GetMatrix(entry)
	matrix.Invert()

	viewWidth := safeArea.Max.X - safeArea.Min.X
	viewHeight := safeArea.Max.Y - safeArea.Min.Y

	// Center and offset to safe area
	matrix.Translate(
		(viewWidth/2)+safeArea.Min.X,
		(viewHeight/2)+safeArea.Min.Y,
	)
	return matrix
}

func (c *camera) ToWorld(ecs donburi.World, screen Screen, screenPos geom.Vec2) geom.Vec2 {
	invView := c.view(ecs, screen)
	invView.Invert()

	worldX, worldY := invView.Apply(screenPos.X, screenPos.Y)
	return geom.Vec2{X: worldX, Y: worldY}
}

func (c *camera) ToScreen(ecs donburi.World, screen Screen, worldPos geom.Vec2) geom.Vec2 {
	view := c.view(ecs, screen)

	screenX, screenY := view.Apply(worldPos.X, worldPos.Y)
	return geom.Vec2{X: screenX, Y: screenY}
}
