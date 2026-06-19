package debug

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DebugWorld interface {
	GetHashCells(area geom.AABB) []geom.AABB
}

var path vector.Path

func DrawWorldGrid(
	screen *ebiten.Image,
	world engine.World,
	camera engine.Camera,
	viewport geom.AABB,
	viewMatrix ebiten.GeoM) {
	debugworld, ok := world.(DebugWorld)
	if !ok {
		return
	}

	path.Reset()

	cells := debugworld.GetHashCells(viewport)
	for _, cell := range cells {
		if !cell.Intersects(viewport) {
			continue
		}

		min := camera.ToScreen(cell.Min)
		max := camera.ToScreen(cell.Max)

		path.MoveTo(float32(min.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(max.Y))
		path.LineTo(float32(min.X), float32(max.Y))
		path.Close()
	}

	vector.StrokePath(screen, &path, &vector.StrokeOptions{Width: 2}, &vector.DrawPathOptions{})
}
