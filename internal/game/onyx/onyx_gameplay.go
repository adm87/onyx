package onyx

import (
	"fmt"
	"image/color"
	"time"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/aseprite"
	asepritemodule "github.com/adm87/onyx/pkg/aseprite"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	tiledmodule "github.com/adm87/onyx/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

var (
	debugEntityInfo   = false
	debugColliderInfo = true
)

var gameplayManifest = []file.FilePath{
	content.AssetsTiledGym04,
	content.AssetsAsepriteCaptainImg,
	content.AssetsAsepriteCaptainJson,
}

func (o *Onyx) GameplayScene() engine.SceneState {
	var tilemapEntry *donburi.Entry
	var spriteEntry *donburi.Entry
	var tilemap *tiledmodule.Tilemap
	var path vector.Path
	var move geom.Vec2
	var tilemapHandle uint64
	var err error
	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			a := o.game.Assets()
			c := o.game.Camera()
			i := o.images
			t := o.tiled
			as := o.aseprite

			o.game.Screen().SetBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255})

			err = a.Load(content.AssetsFS(), gameplayManifest...)
			assert.Nil(err, fmt.Sprintf("failed to load gameplay assets: %v", err))

			tmxHandle, exists := t.GetTmxHandle(content.AssetsTiledGym04)
			assert.True(exists, "failed to get handle for tiled map")

			tilemap, tilemapHandle = t.BuildTilemap(tmxHandle)
			tilemapEntry = t.CreateTilemapEntity(ecs, tiledmodule.WithTilemapHandle(tilemapHandle))

			playerImgHandle, exists := i.GetAssetHandle(content.AssetsAsepriteCaptainImg)
			assert.True(exists, "failed to get handle for player image")

			playerDataHandle, exists := a.GetDataHandle(content.AssetsAsepriteCaptainJson)
			assert.True(exists, "failed to get handle for player data")

			playerData, exists := a.GetData(playerDataHandle)
			assert.True(exists, "failed to get data for player")

			as.BuildAnimations(playerImgHandle, playerData)

			spriteEntry = as.CreateSpriteEntity(ecs,
				asepritemodule.WithImageHandle(playerImgHandle),
				asepritemodule.WithClip("Run"),
				asepritemodule.Playing(),
			)
			images.SetAnchor(spriteEntry, 0.5, 1.0)
			rendering.SetLayer(spriteEntry, 6)

			playerWidth, playerHeight, _ := i.GetFrameSize(playerImgHandle, 0)
			transform.SetLocalBounds(spriteEntry, &geom.AABB{
				Min: geom.Vec2{X: -float64(playerWidth) / 2, Y: -float64(playerHeight)},
				Max: geom.Vec2{X: float64(playerWidth) / 2, Y: 0},
			})

			transform.SetPosition(spriteEntry, tilemap.Bounds().Center())

			c.SetPosition(tilemap.Bounds().Center())
			c.SetZoom(0.25)

			world := o.game.World()
			world.Add(tilemapEntry)
			world.Add(spriteEntry)

			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
				return engine.SceneExitNone, ebiten.Termination
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF) {
				ebiten.SetFullscreen(!ebiten.IsFullscreen())
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
				debugEntityInfo = !debugEntityInfo
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
				debugColliderInfo = !debugColliderInfo
			}

			move = geom.Vec2{}

			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				move.X -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				move.X += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				move.Y -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) {
				move.Y += 1
			}

			if move.X != 0 || move.Y != 0 {
				aseprite.SetClip(spriteEntry, "Run")
				if move.X < 0 {
					transform.SetScale(spriteEntry, -1, 1)
				} else {
					transform.SetScale(spriteEntry, 1, 1)
				}
			} else {
				aseprite.SetClip(spriteEntry, "Idle")
			}
			return engine.SceneExitNone, nil
		},
		OnFixedUpdate: func(ecs donburi.World, dt float64) error {
			if move.X != 0 || move.Y != 0 {
				position := transform.GetPosition(spriteEntry)
				move = move.Normalize().Mul(100 * dt)

				newPos := position.Add(move)
				transform.SetPosition(spriteEntry, newPos)

				o.game.World().Update(spriteEntry)
			}
			return nil
		},
		OnLateUpdate: func(ecs donburi.World, dt float64) error {
			d := time.Duration(float64(time.Second) * dt)
			o.game.World().QueryRegion(ecs, o.game.Camera().Viewport(), func(entry *donburi.Entry) {
				if aseprite.IsPlaying(entry) {
					o.aseprite.UpdateAnimation(entry, d)
				}
			})
			return nil
		},
		OnRender: func(entries []*donburi.Entry, img *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error {
			camera := o.game.Camera()

			if debugEntityInfo {
				path.Reset()
				drawDebugInfo(&path, camera, entries, img)
			}

			min := o.game.Screen().SafeArea().Min
			ebitenutil.DebugPrintAt(img, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), int(min.X), int(min.Y))
			return nil
		},
	}
}

func drawDebugInfo(path *vector.Path, camera engine.Camera, entries []*donburi.Entry, img *ebiten.Image) {
	for _, entry := range entries {
		trans := transform.GetTransform(entry)
		aabb := transform.GetLocalBounds(entry).Translate(trans.X, trans.Y)

		min := camera.ToScreen(aabb.Min)
		max := camera.ToScreen(aabb.Max)

		path.MoveTo(float32(min.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(min.Y))
		path.LineTo(float32(max.X), float32(max.Y))
		path.LineTo(float32(min.X), float32(max.Y))
		path.Close()

		ebitenutil.DebugPrintAt(img, fmt.Sprintf("Entity:\n  ID %d", entry.Entity()), int(min.X), int(min.Y))
	}
	vector.StrokePath(img, path, &vector.StrokeOptions{Width: 2}, &vector.DrawPathOptions{})
}
