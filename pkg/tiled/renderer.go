package tiled

import (
	"image"
	"math"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
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

	buffers map[engine.FilePath]*ebiten.Image
	drawn   map[engine.FilePath]struct{}
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
		buffers:           make(map[engine.FilePath]*ebiten.Image),
		drawn:             make(map[engine.FilePath]struct{}),
	}
}

func (a *TiledRenderingAdapter) getBuffer(ref engine.FilePath) *ebiten.Image {
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

func (a *TiledRenderingAdapter) GetRenderTasks(world donburi.World, viewMatrix ebiten.GeoM) []engine.RenderTask {
	a.renderingTasks = a.renderingTasks[:0]

	// Transform screen corners to world coordinates for culling
	worldMin := a.camera.ToWorld(world, a.screen, a.screen.SafeArea().Min)
	worldMax := a.camera.ToWorld(world, a.screen, a.screen.SafeArea().Max)

	clear(a.drawn)

	// Iterate over all entities with a Tilemap component and enqueue render tasks for visible tiles
	TilemapQuery.Each(world, func(entry *donburi.Entry) {
		ref := GetTilemapRef(entry)

		visible := rendering.IsVisible(entry)
		if !visible {
			return // Skip invisible tilemaps
		}

		filter := rendering.GetFilter(entry)
		layer := rendering.GetLayer(entry)
		zIndex := rendering.GetZIndex(entry)

		// Get the tilemap buffer and clear it, this will ensure we don't have artifacts from previous frames when rendering the tilemap
		buffer := a.getBuffer(ref)
		buffer.Clear()

		// Mark the tilemap as drawn before desiding if it's visible.
		// We do this to ensure the buffer for a valid tilemap doesn't get deallocated in case it becomes viable in the next frame.
		a.drawn[ref] = struct{}{}

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
	})

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
			tsxPath := engine.FilePath(tileset.Source)

			tsx, exists := a.tiledAssetAdapter.tsxCache[tsxPath]
			if !exists {
				continue // Skip tiles that reference a tileset without tsx data
			}

			tilesetImg, exists := a.imageAssetAdapter.GetImage(engine.FilePath(tsx.Image.Source))
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
