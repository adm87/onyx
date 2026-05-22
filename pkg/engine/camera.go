package engine

import (
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var CameraTag = donburi.NewTag("Camera")

type Camera interface {
	Position() geom.Vec2
	SetPosition(pos geom.Vec2)
	ToWorld(screen geom.Vec2) geom.Vec2
	ToScreen(world geom.Vec2) geom.Vec2
}

type camera struct {
	world  donburi.World
	entity donburi.Entity
	screen *screen
}

func newCamera(world donburi.World, screen *screen) *camera {
	entity := world.Create(
		transform.Position,
		transform.Rotation,
		transform.Scale,
		transform.Matrix,
		CameraTag,
	)
	return &camera{
		world:  world,
		entity: entity,
		screen: screen,
	}
}

func (c *camera) Position() geom.Vec2 {
	entry := c.world.Entry(c.entity)
	return transform.GetPosition(entry)
}

func (c *camera) SetPosition(pos geom.Vec2) {
	entry := c.world.Entry(c.entity)
	transform.SetPosition(entry, pos)
}

func (c *camera) view() ebiten.GeoM {
	entry := c.world.Entry(c.entity)

	matrix := transform.GetMatrix(entry)
	matrix.Invert()

	viewWidth := c.screen.safeArea.Max.X - c.screen.safeArea.Min.X
	viewHeight := c.screen.safeArea.Max.Y - c.screen.safeArea.Min.Y

	matrix.Translate(viewWidth/2, viewHeight/2)
	matrix.Translate(c.screen.safeArea.Min.X, c.screen.safeArea.Min.Y)

	return matrix
}

func (c *camera) ToWorld(screenPos geom.Vec2) geom.Vec2 {
	invView := c.view()
	invView.Invert()

	worldX, worldY := invView.Apply(screenPos.X, screenPos.Y)
	return geom.Vec2{X: worldX, Y: worldY}
}

func (c *camera) ToScreen(worldPos geom.Vec2) geom.Vec2 {
	view := c.view()

	screenX, screenY := view.Apply(worldPos.X, worldPos.Y)
	return geom.Vec2{X: screenX, Y: screenY}
}
