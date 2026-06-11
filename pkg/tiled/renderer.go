package tiled

import (
	"image"
	"math"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type renderingAdapter struct {
	screen         engine.Screen
	imageModule    *images.ImageModule
	tiledAssets    *assetAdapter
	renderingTasks []engine.RenderingTask
	buffers        map[uint64][]*ebiten.Image
}

func newRenderingAdapter(
	screen engine.Screen,
	imageModule *images.ImageModule,
	tiledAssets *assetAdapter) *renderingAdapter {
	return &renderingAdapter{
		screen:         screen,
		imageModule:    imageModule,
		tiledAssets:    tiledAssets,
		renderingTasks: make([]engine.RenderingTask, 1),
		buffers:        make(map[uint64][]*ebiten.Image),
	}
}

func (a *renderingAdapter) getBuffer(handle uint64, layer int) *ebiten.Image {
	width, height := a.screen.SafeArea().Width(), a.screen.SafeArea().Height()

	buffers, exists := a.buffers[handle]
	if !exists {
		buffers = make([]*ebiten.Image, 0, 4)
	}

	for len(buffers) <= layer {
		buffers = append(buffers, ebiten.NewImage(int(width), int(height)))
	}

	a.buffers[handle] = buffers
	buffer := buffers[layer]

	if buffer.Bounds().Dx() == int(width) && buffer.Bounds().Dy() == int(height) {
		return buffer
	}

	buffer.Deallocate()
	buffer = ebiten.NewImage(int(width), int(height))
	buffers[layer] = buffer

	return buffer
}

func (a *renderingAdapter) releaseBuffer(handle uint64) {
	if buffers, exists := a.buffers[handle]; exists {
		for _, buffer := range buffers {
			buffer.Deallocate()
		}
		delete(a.buffers, handle)
	}
}

func (a *renderingAdapter) GetRenderingTasks(entry *donburi.Entry, viewport geom.AABB, viewMatrix ebiten.GeoM) []engine.RenderingTask {
	a.renderingTasks = a.renderingTasks[:0]

	handle, exists := GetTilemapHandle(entry)
	if !exists {
		return a.renderingTasks
	}

	layer := rendering.GetLayer(entry)
	zindex := rendering.GetZIndex(entry)
	filter := rendering.GetFilter(entry)
	color := rendering.GetColor(entry)

	tilemap, exists := a.tiledAssets.tilemapStore.Get(handle)
	if !exists {
		return a.renderingTasks
	}

	tmx, exists := a.tiledAssets.tmxStore.Get(handle)
	if !exists {
		return a.renderingTasks
	}

	var buffer *ebiten.Image

	minTileX := int(math.Floor(viewport.Min.X / float64(tmx.TileWidth)))
	maxTileX := int(math.Floor(viewport.Max.X / float64(tmx.TileWidth)))
	minTileY := int(math.Floor(viewport.Min.Y / float64(tmx.TileHeight)))
	maxTileY := int(math.Floor(viewport.Max.Y / float64(tmx.TileHeight)))

	for i := range tilemap.layers {
		if !tmx.Layers[i].Visible {
			continue
		}

		buffer = a.getBuffer(handle, i)
		buffer.Clear()

		a.drawTilemapLayer(
			buffer,
			tilemap, i,
			tmx.TileWidth, tmx.TileHeight,
			minTileX, maxTileX,
			minTileY, maxTileY,
			tmx.Tilesets,
			viewMatrix)
		a.renderingTasks = append(a.renderingTasks, engine.RenderingTask{
			Layer:  layer,
			ZIndex: zindex + i,
			Job: func(target *ebiten.Image) {
				opt := &ebiten.DrawImageOptions{
					Filter: filter,
				}

				opt.ColorScale.ScaleWithColor(color)
				opt.ColorScale.ScaleAlpha(float32(color.A) / 255)

				target.DrawImage(buffer, opt)
			},
		})
	}

	return a.renderingTasks
}

func (a *renderingAdapter) drawTilemapLayer(
	target *ebiten.Image,
	tilemap *Tilemap,
	layerIndex int,
	cellWidth, cellHeight int,
	minTileX, maxTileX int,
	minTileY, maxTileY int,
	tilesets []TmxTileset,
	viewMatrix ebiten.GeoM) {

	opts := &ebiten.DrawImageOptions{}

	for y := minTileY; y <= maxTileY; y++ {
		for x := minTileX; x <= maxTileX; x++ {
			tile, exists := tilemap.GetTile(layerIndex, x, y)
			if !exists || tile.ID() == 0 {
				continue
			}

			tileset := NearestTileset(tilesets, tile.ID())
			tsx, exists := a.tiledAssets.tsxStore.Get(tileset.Handle)
			if !exists {
				continue
			}

			img, exists := a.imageModule.GetImage(tsx.Image.Handle)
			if !exists {
				continue
			}

			tileID := tile.ID() - uint32(tileset.FirstGID)
			tileX, tileY := x*cellWidth, y*cellHeight

			srcX := int(tileID % uint32(tsx.Columns) * uint32(tsx.TileWidth))
			srcY := int(tileID / uint32(tsx.Columns) * uint32(tsx.TileHeight))

			opts.GeoM.Reset()

			// Ref: https://doc.mapeditor.org/en/stable/reference/global-tile-ids/#tile-flipping
			if tile.FlippedDiagonally() {
				opts.GeoM.Rotate(math.Pi * 0.5)
				opts.GeoM.Scale(-1, 1)
				opts.GeoM.Translate(float64(tsx.TileHeight-tsx.TileWidth), 0)
			}
			if tile.FlippedHorizontally() {
				opts.GeoM.Scale(-1, 1)
				opts.GeoM.Translate(float64(tsx.TileWidth), 0)
			}
			if tile.FlippedVertically() {
				opts.GeoM.Scale(1, -1)
				opts.GeoM.Translate(0, float64(tsx.TileHeight))
			}

			opts.GeoM.Translate(0, float64(cellHeight-tsx.TileHeight))
			opts.GeoM.Translate(float64(tsx.TileOffset.X), float64(tsx.TileOffset.Y))
			opts.GeoM.Translate(float64(tileX), float64(tileY))
			opts.GeoM.Concat(viewMatrix)

			target.DrawImage(img.SubImage(
				image.Rect(
					srcX, srcY,
					srcX+tsx.TileWidth,
					srcY+tsx.TileHeight,
				),
			).(*ebiten.Image), opts)
		}
	}
}
