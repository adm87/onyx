package tiled

import (
	"image"
	"image/color"
	"math"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/asset"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type TiledRenderingAdapter struct {
	tiledAssetAdapter *TiledAssetAdapter
	imageAssetAdapter *images.ImageAssetAdapter

	camera engine.Camera
	screen engine.Screen

	renderingTasks []engine.RenderTask

	buffers map[file.FilePath]*ebiten.Image
	drawn   map[file.FilePath]struct{}
}

func NewTiledRenderingAdapter(
	tiledAssetAdapter *TiledAssetAdapter,
	imageAssetAdapter *images.ImageAssetAdapter,
	camera engine.Camera,
	screen engine.Screen) *TiledRenderingAdapter {
	return &TiledRenderingAdapter{
		tiledAssetAdapter: tiledAssetAdapter,
		imageAssetAdapter: imageAssetAdapter,
		camera:            camera,
		screen:            screen,
		renderingTasks:    make([]engine.RenderTask, 0, 10),
		buffers:           make(map[file.FilePath]*ebiten.Image),
		drawn:             make(map[file.FilePath]struct{}),
	}
}

func (a *TiledRenderingAdapter) getBuffer(ref file.FilePath) *ebiten.Image {
	screenSize := a.screen.Size()
	buffer, exists := a.buffers[ref]

	if exists {
		bufWidth, bufHeight := buffer.Bounds().Dx(), buffer.Bounds().Dy()
		if bufWidth == int(screenSize.X) && bufHeight == int(screenSize.Y) {
			return buffer
		}
		buffer.Deallocate()
	}

	buffer = ebiten.NewImage(
		int(screenSize.X),
		int(screenSize.Y),
	)

	a.buffers[ref] = buffer
	return buffer
}

func (a *TiledRenderingAdapter) GetRenderTasks(ecs donburi.World, viewMatrix ebiten.GeoM) []engine.RenderTask {
	a.renderingTasks = a.renderingTasks[:0]

	// Transform screen corners to world coordinates for culling
	worldMin := a.camera.ToWorld(ecs, a.screen, a.screen.SafeArea().Min)
	worldMax := a.camera.ToWorld(ecs, a.screen, a.screen.SafeArea().Max)

	clear(a.drawn)

	rendering.QueryWith(ecs, TiledQuery,
		func(entry *donburi.Entry, anchor geom.Vec2, color color.RGBA, filter ebiten.Filter, visible bool, layer, zIndex int) {
			// TODO - apply color tint, anchor, and filter to tilemap rendering (currently ignored)

			ref := asset.GetAssetReference(entry)
			if ref == asset.UnknownRef {
				return // Don't enqueue render tasks for entities without a tilemap reference
			}

			// Clear the buffer and mark it as drawn for this frame.
			// Buffers that are not marked as drawn after processing all visible tilemaps will be deallocated to free up memory.
			buffer := a.getBuffer(ref)
			buffer.Clear()

			// To avoid reallocating buffers every frame for invisible tilemaps,
			// we still need to mark them as drawn even if they are not visible.
			// This way, they will be reused when they become visible again instead of being deallocated and reallocated.
			// Only when a tilemap entity is removed or no longer references the same tilemap will its buffer be deallocated.
			a.drawn[ref] = struct{}{}

			if !visible {
				return // Skip invisible tilemaps
			}

			tilemap, exists := a.tiledAssetAdapter.tilemaps[ref]
			if !exists {
				return // Don't enqueue render tasks for entities with an invalid tilemap reference
			}

			tmx, exists := a.tiledAssetAdapter.tmxCache[ref]
			if !exists {
				return // Don't enqueue render tasks for entities with an invalid tilemap reference
			}

			minTileX := int(math.Floor(worldMin.X / float64(tmx.TileWidth)))
			maxTileX := int(math.Floor(worldMax.X / float64(tmx.TileWidth)))
			minTileY := int(math.Floor(worldMin.Y / float64(tmx.TileHeight)))
			maxTileY := int(math.Floor(worldMax.Y / float64(tmx.TileHeight)))

			for i := range tilemap.layers {
				if !tmx.Layers[i].Visible {
					continue // Skip invisible layers
				}
				a.drawLayerToBuffer(
					buffer,
					tilemap,
					tmx.Tilesets,
					i,
					tmx.TileWidth, tmx.TileHeight,
					minTileX, maxTileX,
					minTileY, maxTileY,
					viewMatrix,
				)
			}

			a.renderingTasks = append(a.renderingTasks, engine.RenderTask{
				Render: func(screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
					screen.DrawImage(buffer, &ebiten.DrawImageOptions{
						Filter: filter,
					})
					return nil
				},
				Layer:  layer,
				ZIndex: zIndex,
			})
		},
	)

	// Deallocate buffers that were detected to be no longer viable.
	for ref, buffer := range a.buffers {
		if _, drawn := a.drawn[ref]; !drawn {
			buffer.Deallocate()
			delete(a.buffers, ref)
		}
	}

	return a.renderingTasks
}

func (a *TiledRenderingAdapter) drawLayerToBuffer(
	buffer *ebiten.Image,
	tilemap *Tilemap,
	tilesets []TmxTileset,
	layer int,
	cellWidth, cellHeight int,
	minTileX, maxTileX int,
	minTileY, maxTileY int,
	viewMatrix ebiten.GeoM,
) {
	opt := ebiten.DrawImageOptions{}

	for y := minTileY; y <= maxTileY; y++ {
		for x := minTileX; x <= maxTileX; x++ {
			tile, ok := tilemap.GetTile(layer, x, y)
			if !ok {
				continue // Skip tiles that are outside the bounds of the tile array
			}

			if tile.ID() == 0 {
				continue // Skip empty tiles
			}

			tileset := NearestTileset(tilesets, tile.ID())
			tsxPath := file.FilePath(tileset.Source)

			tsx, exists := a.tiledAssetAdapter.tsxCache[tsxPath]
			if !exists {
				continue // Skip tiles that reference a tileset without tsx data
			}

			tilesetImg, exists := a.imageAssetAdapter.GetImage(file.FilePath(tsx.Image.Source))
			if !exists {
				continue // Skip tiles that reference a missing tileset image
			}

			tileX, tileY := x*cellWidth, y*cellHeight
			tileID := tile.ID() - uint32(tileset.FirstGID)

			srcX := int(tileID % uint32(tsx.Columns) * uint32(tsx.TileWidth))
			srcY := int(tileID / uint32(tsx.Columns) * uint32(tsx.TileHeight))

			opt.GeoM.Reset()

			// Ref: https://doc.mapeditor.org/en/stable/reference/global-tile-ids/#tile-flipping
			if tile.FlippedDiagonally() {
				opt.GeoM.Rotate(math.Pi * 0.5)
				opt.GeoM.Scale(-1, 1)
				opt.GeoM.Translate(float64(tsx.TileHeight-tsx.TileWidth), 0)
			}
			if tile.FlippedHorizontally() {
				opt.GeoM.Scale(-1, 1)
				opt.GeoM.Translate(float64(tsx.TileWidth), 0)
			}
			if tile.FlippedVertically() {
				opt.GeoM.Scale(1, -1)
				opt.GeoM.Translate(0, float64(tsx.TileHeight))
			}

			opt.GeoM.Translate(0, float64(cellHeight-tsx.TileHeight))
			opt.GeoM.Translate(float64(tsx.TileOffset.X), float64(tsx.TileOffset.Y))
			opt.GeoM.Translate(float64(tileX), float64(tileY))
			opt.GeoM.Concat(viewMatrix)

			buffer.DrawImage(tilesetImg.SubImage(
				image.Rect(
					srcX, srcY,
					srcX+tsx.TileWidth,
					srcY+tsx.TileHeight,
				),
			).(*ebiten.Image), &opt)
		}
	}
}
