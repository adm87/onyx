package tiled

import (
	"image"
	"math"

	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/adm87/onyx-game/pkg/engine/components/rendering"
	"github.com/adm87/onyx-game/pkg/images"
	"github.com/adm87/onyx-game/pkg/tiled/components"
	"github.com/adm87/onyx-game/pkg/tiled/data"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type TiledRenderingAdapter struct {
	tiledAssetAdapter *TiledAssetAdapter
	imageAssetAdapter *images.ImageAssetAdapter

	screen engine.Screen
	logger engine.Logger

	renderingTasks []engine.RenderTask

	buffers map[engine.FilePath]*ebiten.Image
	drawn   map[engine.FilePath]struct{}
}

func NewTiledRenderingAdapter(
	tiledAssetAdapter *TiledAssetAdapter,
	imageAssetAdapter *images.ImageAssetAdapter,
	screen engine.Screen,
	logger engine.Logger) *TiledRenderingAdapter {
	return &TiledRenderingAdapter{
		tiledAssetAdapter: tiledAssetAdapter,
		imageAssetAdapter: imageAssetAdapter,
		screen:            screen,
		logger:            logger,
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

	// Invert the view matrix to transform screen coordinates back to world coordinates for culling
	invViewMatrix := viewMatrix
	invViewMatrix.Invert()

	// Calculate the world coordinates of the corners of the screen to determine which tiles are visible
	screenMinX, screenMinY := a.screen.SafeArea().Min.XY()
	screenMaxX, screenMaxY := a.screen.SafeArea().Max.XY()

	// Transform screen corners to world coordinates for culling
	worldMinX, worldMinY := invViewMatrix.Apply(screenMinX, screenMinY)
	worldMaxX, worldMaxY := invViewMatrix.Apply(screenMaxX, screenMaxY)

	clear(a.drawn)

	// Iterate over all entities with a Tilemap component and enqueue render tasks for visible tiles
	components.TilemapQuery.Each(world, func(entry *donburi.Entry) {
		ref := components.GetTilemapRef(entry)

		tilemap, exists := a.tiledAssetAdapter.tilemaps[ref]
		if !exists {
			a.logger.Error("tilemap asset not found for reference: %s", ref)
			return // Don't enqueue render tasks for entities with an invalid tilemap reference
		}

		tmx, exists := a.tiledAssetAdapter.tmxCache[ref]
		if !exists {
			a.logger.Error("tmx data not found for tilemap reference: %s", ref)
			return // Don't enqueue render tasks for entities with an invalid tilemap reference
		}

		minTileX := int(math.Floor(worldMinX / float64(tmx.TileWidth)))
		maxTileX := int(math.Ceil(worldMaxX / float64(tmx.TileWidth)))
		minTileY := int(math.Floor(worldMinY / float64(tmx.TileHeight)))
		maxTileY := int(math.Ceil(worldMaxY / float64(tmx.TileHeight)))

		if minTileX > int(tilemap.tileBounds.Max.X) || maxTileX < int(tilemap.tileBounds.Min.X) ||
			minTileY > int(tilemap.tileBounds.Max.Y) || maxTileY < int(tilemap.tileBounds.Min.Y) {
			return // Don't enqueue render tasks for tilemaps that are completely outside the view
		}

		filter := rendering.GetFilter(entry)
		layer := rendering.GetLayer(entry)
		zIndex := rendering.GetZIndex(entry)

		buffer := a.getBuffer(ref)
		buffer.Clear()

		for i := range tilemap.layers {
			if !tmx.Layers[i].Visible {
				continue // Skip invisible layers
			}
			a.drawLayerToBuffer(
				buffer,
				tilemap,
				tmx,
				i,
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

		a.drawn[ref] = struct{}{}
	})
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
	tmx *data.Tmx,
	layer int,
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
			if tile.id == 0 {
				continue // Skip empty tiles
			}

			tileset := data.NearestTileset(tmx.Tilesets, tile.id)
			tsxPath := engine.FilePath(tileset.Source)

			tsx, exists := a.tiledAssetAdapter.tsxCache[tsxPath]
			if !exists {
				a.logger.Error("tsx data not found for tileset image: %s", tileset.Source)
				continue // Skip tiles that reference a tileset without tsx data
			}

			tilesetImg, exists := a.imageAssetAdapter.GetImage(engine.FilePath(tsx.Image.Source))
			if !exists {
				a.logger.Error("failed to get image asset for tileset: %s", tileset.Source)
				continue // Skip tiles that reference a missing tileset image
			}

			tileX, tileY := x*tsx.TileWidth, y*tsx.TileHeight
			tileID := tile.id - uint32(tileset.FirstGID)

			srcX := int(tileID % uint32(tsx.Columns) * uint32(tsx.TileWidth))
			srcY := int(tileID / uint32(tsx.Columns) * uint32(tsx.TileHeight))

			opt.GeoM.Reset()
			opt.GeoM.Translate(float64(tileX), float64(tileY))
			opt.GeoM.Concat(viewMatrix)

			buffer.DrawImage(tilesetImg.SubImage(
				image.Rect(
					srcX, srcY,
					srcX+tmx.TileWidth,
					srcY+tmx.TileHeight,
				),
			).(*ebiten.Image), &opt)
		}
	}
}
