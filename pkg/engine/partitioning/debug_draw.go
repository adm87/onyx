package partitioning

import (
	"image/color"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawDebug[T comparable](screen *ebiten.Image, hash *SpatialHash[T], safeArea geom.AABB, viewMatrix ebiten.GeoM) {
	ebitenutil.DebugPrintAt(screen, "Spatial Hash Debug", int(safeArea.Min.X), int(safeArea.Min.Y))

	invViewMatrix := viewMatrix
	invViewMatrix.Invert()

	worldMinX, worldMinY := invViewMatrix.Apply(safeArea.Min.X, safeArea.Min.Y)
	worldMaxX, worldMaxY := invViewMatrix.Apply(safeArea.Max.X, safeArea.Max.Y)

	for _, grid := range hash.grids {
		cellSize := grid.cellSize

		startCellX := int64(worldMinX / cellSize)
		startCellY := int64(worldMinY / cellSize)
		endCellX := int64(worldMaxX / cellSize)
		endCellY := int64(worldMaxY / cellSize)

		for x := startCellX; x <= endCellX; x++ {
			for y := startCellY; y <= endCellY; y++ {
				cellCoord := encodeCoord(x, y)
				if _, exists := grid.cells[cellCoord]; exists {
					cellMinX := float64(x) * cellSize
					cellMinY := float64(y) * cellSize
					cellMaxX := cellMinX + cellSize
					cellMaxY := cellMinY + cellSize

					screenMinX, screenMinY := viewMatrix.Apply(cellMinX, cellMinY)
					screenMaxX, screenMaxY := viewMatrix.Apply(cellMaxX, cellMaxY)

					vector.StrokeRect(screen,
						float32(screenMinX),
						float32(screenMinY),
						float32(screenMaxX-screenMinX),
						float32(screenMaxY-screenMinY),
						2,
						color.RGBA{255, 0, 0, 100},
						false,
					)
				}
			}
		}
	}
}
