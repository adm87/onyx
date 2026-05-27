package tiled

import (
	"math"

	"github.com/adm87/onyx-game/pkg/engine/geom"
	"github.com/adm87/onyx-game/pkg/tiled/data"
)

type Tile struct {
	id    uint32
	flags uint32
}

func (t Tile) FlippedHorizontally() bool {
	return t.flags&FlippedHorizontallyFlag != 0
}

func (t Tile) FlippedVertically() bool {
	return t.flags&FlippedVerticallyFlag != 0
}

func (t Tile) FlippedDiagonally() bool {
	return t.flags&FlippedDiagonallyFlag != 0
}

func (t Tile) RotatedHexagonal120() bool {
	return t.flags&RotatedHexagonal120Flag != 0
}

type Tilemap struct {
	bounds geom.AABB
	layers []TilemapLayer
}

type TilemapLayer struct {
	name   string
	bounds geom.AABB
	tiles  []Tile
}

func buildTilemap(tmx *data.Tmx) (*Tilemap, error) {
	min, max := findMapSize(tmx)

	tilemap := &Tilemap{
		bounds: geom.AABB{
			Min: min,
			Max: max,
		},
		layers: make([]TilemapLayer, len(tmx.Layers)),
	}

	for i, layer := range tmx.Layers {
		var newLayer TilemapLayer

		if err := buildTilemapLayer(layer, &newLayer); err != nil {
			return nil, err
		}

		tilemap.layers[i] = newLayer
	}

	return tilemap, nil
}

func findMapSize(tmx *data.Tmx) (geom.Vec2, geom.Vec2) {
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := -math.MaxInt, -math.MaxInt

	for _, layer := range tmx.Layers {
		layerMin, layerMax := findLayerSize(layer)

		if layerMin.X < float64(minX) {
			minX = int(layerMin.X)
		}
		if layerMin.Y < float64(minY) {
			minY = int(layerMin.Y)
		}
		if layerMax.X > float64(maxX) {
			maxX = int(layerMax.X)
		}
		if layerMax.Y > float64(maxY) {
			maxY = int(layerMax.Y)
		}
	}

	return geom.Vec2{
			X: float64(minX),
			Y: float64(minY),
		}, geom.Vec2{
			X: float64(maxX),
			Y: float64(maxY),
		}
}

func findLayerSize(layer data.TmxLayer) (geom.Vec2, geom.Vec2) {
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := -math.MaxInt, -math.MaxInt

	if len(layer.Data.Chunks) == 0 {
		minX, minY = 0, 0
		maxX, maxY = layer.Width, layer.Height
	} else {
		for _, chunk := range layer.Data.Chunks {
			if chunk.X < minX {
				minX = chunk.X
			}
			if chunk.Y < minY {
				minY = chunk.Y
			}
			if chunk.X+chunk.Width > maxX {
				maxX = chunk.X + chunk.Width
			}
			if chunk.Y+chunk.Height > maxY {
				maxY = chunk.Y + chunk.Height
			}
		}
	}

	return geom.Vec2{
			X: float64(minX),
			Y: float64(minY),
		}, geom.Vec2{
			X: float64(maxX),
			Y: float64(maxY),
		}
}

func buildTilemapLayer(tmxLayer data.TmxLayer, layer *TilemapLayer) error {
	min, max := findLayerSize(tmxLayer)

	layer.name = tmxLayer.Name
	layer.bounds = geom.AABB{
		Min: min,
		Max: max,
	}

	tiles, err := decodeContent(
		tmxLayer.Data.Encoding,
		tmxLayer.Data.Compression,
		tmxLayer.Data.Content,
	)
	if err != nil {
		return err
	}
	layer.tiles = tiles

	return nil
}
