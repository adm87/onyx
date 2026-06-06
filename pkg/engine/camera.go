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
	Position() geom.Vec2
	SetPosition(pos geom.Vec2)
	Zoom() float64
	SetZoom(zoom float64)
	ToWorld(screen Screen, position geom.Vec2) geom.Vec2
	ToScreen(screen Screen, position geom.Vec2) geom.Vec2
}

type camera struct {
	donburi.Entity

	world  *world
	screen *screen
}

func newCamera(world *world, screen *screen, entity donburi.Entity) *camera {
	return &camera{
		Entity: entity,
		world:  world,
		screen: screen,
	}
}

func (c *camera) Position() geom.Vec2 {
	entry := c.world.ecs.Entry(c.Entity)
	return transform.GetPosition(entry)
}

func (c *camera) SetPosition(pos geom.Vec2) {
	entry := c.world.ecs.Entry(c.Entity)
	transform.SetPosition(entry, pos)
}

func (c *camera) Zoom() float64 {
	entry := c.world.ecs.Entry(c.Entity)
	scale := transform.GetScale(entry)
	return scale.X
}

func (c *camera) SetZoom(zoom float64) {
	entry := c.world.ecs.Entry(c.Entity)
	transform.SetScale(entry, geom.Vec2{X: zoom, Y: zoom})
}

func (c *camera) view(min, max geom.Vec2) ebiten.GeoM {
	entry := c.world.ecs.Entry(c.Entity)

	matrix := transform.GetMatrix(entry)
	matrix.Invert()

	viewWidth := max.X - min.X
	viewHeight := max.Y - min.Y

	// Center and offset to safe area
	matrix.Translate(
		(viewWidth/2)+min.X,
		(viewHeight/2)+min.Y,
	)
	return matrix
}

func (c *camera) ToWorld(screen Screen, screenPos geom.Vec2) geom.Vec2 {
	invView := c.view(
		screen.SafeArea().Min,
		screen.SafeArea().Max,
	)
	invView.Invert()

	worldX, worldY := invView.Apply(screenPos.X, screenPos.Y)
	return geom.Vec2{X: worldX, Y: worldY}
}

func (c *camera) ToScreen(screen Screen, worldPos geom.Vec2) geom.Vec2 {
	view := c.view(
		screen.SafeArea().Min,
		screen.SafeArea().Max,
	)
	screenX, screenY := view.Apply(worldPos.X, worldPos.Y)
	return geom.Vec2{X: screenX, Y: screenY}
}
