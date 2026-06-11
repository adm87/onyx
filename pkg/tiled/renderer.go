package tiled

import (
	"math"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type renderingAdapter struct {
	screen             engine.Screen
	imageModule        *images.ImageModule
	tiledAssets        *assetAdapter
	drawOpts           ebiten.DrawImageOptions
	buffers            map[uint64][]*ebiten.Image
	jobs               []*engine.RenderingJob
	lastMinX, lastMaxX int
	lastMinY, lastMaxY int
}

func newRenderingAdapter(
	screen engine.Screen,
	imageModule *images.ImageModule,
	tiledAssets *assetAdapter) *renderingAdapter {
	return &renderingAdapter{
		screen:      screen,
		imageModule: imageModule,
		tiledAssets: tiledAssets,
		jobs:        make([]*engine.RenderingJob, 0, 100),
		buffers:     make(map[uint64][]*ebiten.Image),
	}
}

func (a *renderingAdapter) getBuffer(handle uint64, layer int) (*ebiten.Image, bool) {
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
		return buffer, false
	}

	buffer.Deallocate()
	buffer = ebiten.NewImage(int(width), int(height))
	buffers[layer] = buffer

	return buffer, true
}

func (a *renderingAdapter) releaseBuffer(handle uint64) {
	if buffers, exists := a.buffers[handle]; exists {
		for _, buffer := range buffers {
			buffer.Deallocate()
		}
		delete(a.buffers, handle)
	}
}

func (a *renderingAdapter) GetJobs(
	entry *donburi.Entry,
	viewport geom.AABB,
	viewMatrix ebiten.GeoM,
	pool engine.RenderingJobPool) []*engine.RenderingJob {
	a.jobs = a.jobs[:0]

	handle, exists := GetTilemapHandle(entry)
	if !exists {
		return a.jobs
	}

	layer := rendering.GetLayer(entry)
	zindex := rendering.GetZIndex(entry)

	tilemap, exists := a.tiledAssets.tilemapStore.Get(handle)
	if !exists {
		return a.jobs
	}

	tmx, exists := a.tiledAssets.tmxStore.Get(handle)
	if !exists {
		return a.jobs
	}

	minTileX := int(math.Floor(viewport.Min.X / float64(tmx.TileWidth)))
	maxTileX := int(math.Floor(viewport.Max.X / float64(tmx.TileWidth)))
	minTileY := int(math.Floor(viewport.Min.Y / float64(tmx.TileHeight)))
	maxTileY := int(math.Floor(viewport.Max.Y / float64(tmx.TileHeight)))

	var buffer *ebiten.Image
	var resized bool

	viewChanged := minTileX != a.lastMinX || maxTileX != a.lastMaxX || minTileY != a.lastMinY || maxTileY != a.lastMaxY
	if viewChanged {
		a.lastMinX, a.lastMaxX, a.lastMinY, a.lastMaxY = minTileX, maxTileX, minTileY, maxTileY
	}

	for i := range tilemap.layers {
		if !tmx.Layers[i].Visible {
			continue
		}

		buffer, resized = a.getBuffer(handle, i)
		if resized || viewChanged {
			buffer.Clear()
			a.drawTilemapLayer(
				buffer,
				tilemap, i,
				tmx.TileWidth, tmx.TileHeight,
				minTileX, maxTileX,
				minTileY, maxTileY,
				tmx.Tilesets,
				viewMatrix,
			)
		}

		job := pool.Get(buffer)
		job.Layer = layer
		job.ZIndex = zindex + i
		a.jobs = append(a.jobs, job)
	}

	return a.jobs
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

	for y := minTileY; y <= maxTileY; y++ {
		for x := minTileX; x <= maxTileX; x++ {
			tile, tileIndex, exists := tilemap.GetTile(layerIndex, x, y)
			if !exists || tile.ID() == 0 {
				continue
			}

			tileset := tilesets[tilemap.tilesets[tileIndex]]

			tsx, exists := a.tiledAssets.tsxStore.Get(tileset.Handle)
			if !exists {
				continue
			}

			a.drawOpts.GeoM.Reset()

			// // Ref: https://doc.mapeditor.org/en/stable/reference/global-tile-ids/#tile-flipping
			if tile.FlippedDiagonally() {
				a.drawOpts.GeoM.Rotate(math.Pi * 0.5)
				a.drawOpts.GeoM.Scale(-1, 1)
				a.drawOpts.GeoM.Translate(float64(tsx.TileHeight-tsx.TileWidth), 0)
			}
			if tile.FlippedHorizontally() {
				a.drawOpts.GeoM.Scale(-1, 1)
				a.drawOpts.GeoM.Translate(float64(tsx.TileWidth), 0)
			}
			if tile.FlippedVertically() {
				a.drawOpts.GeoM.Scale(1, -1)
				a.drawOpts.GeoM.Translate(0, float64(tsx.TileHeight))
			}

			tileX, tileY := x*cellWidth, y*cellHeight
			a.drawOpts.GeoM.Translate(0, float64(cellHeight-tsx.TileHeight))
			a.drawOpts.GeoM.Translate(float64(tsx.TileOffset.X), float64(tsx.TileOffset.Y))
			a.drawOpts.GeoM.Translate(float64(tileX), float64(tileY))
			a.drawOpts.GeoM.Concat(viewMatrix)

			tileID := tile.ID() - uint32(tileset.FirstGID)

			frame, exists := a.imageModule.GetFrame(tsx.Image.Handle, int(tileID))
			if !exists {
				continue
			}

			target.DrawImage(frame, &a.drawOpts)
		}
	}
}
