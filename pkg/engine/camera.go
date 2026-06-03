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
	Position(world donburi.World) geom.Vec2
	SetPosition(world donburi.World, pos geom.Vec2)
	Zoom(world donburi.World) float64
	SetZoom(world donburi.World, zoom float64)
	ToWorld(world donburi.World, screen Screen, position geom.Vec2) geom.Vec2
	ToScreen(world donburi.World, screen Screen, position geom.Vec2) geom.Vec2
}

type camera struct {
	entity donburi.Entity
}

func newCamera(world donburi.World) *camera {
	return &camera{
		entity: world.Create(
			transform.Position,
			transform.Scale,
			transform.Matrix,
			CameraTag,
		),
	}
}

func (c *camera) Position(world donburi.World) geom.Vec2 {
	entry := world.Entry(c.entity)
	return transform.GetPosition(entry)
}

func (c *camera) SetPosition(world donburi.World, pos geom.Vec2) {
	entry := world.Entry(c.entity)
	transform.SetPosition(entry, pos)
}

func (c *camera) Zoom(world donburi.World) float64 {
	entry := world.Entry(c.entity)
	scale := transform.GetScale(entry)
	return scale.X
}

func (c *camera) SetZoom(world donburi.World, zoom float64) {
	entry := world.Entry(c.entity)
	transform.SetScale(entry, geom.Vec2{X: zoom, Y: zoom})
}

func (c *camera) view(world donburi.World, screen Screen) ebiten.GeoM {
	entry := world.Entry(c.entity)
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

func (c *camera) ToWorld(world donburi.World, screen Screen, screenPos geom.Vec2) geom.Vec2 {
	invView := c.view(world, screen)
	invView.Invert()

	worldX, worldY := invView.Apply(screenPos.X, screenPos.Y)
	return geom.Vec2{X: worldX, Y: worldY}
}

func (c *camera) ToScreen(world donburi.World, screen Screen, worldPos geom.Vec2) geom.Vec2 {
	view := c.view(world, screen)

	screenX, screenY := view.Apply(worldPos.X, worldPos.Y)
	return geom.Vec2{X: screenX, Y: screenY}
}
