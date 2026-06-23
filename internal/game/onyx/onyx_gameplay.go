package onyx

import (
	"image/color"
	"time"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/debug"
	"github.com/adm87/onyx/pkg/ecs/camera"
	"github.com/adm87/onyx/pkg/ecs/renderer"
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/aseprite"
	"github.com/adm87/onyx/pkg/plugins/collision"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/adm87/onyx/pkg/plugins/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var gameplayManifest = []file.FilePath{
	content.AssetsAsepriteCaptainImg,
	content.AssetsAsepriteCaptainJson,
	content.AssetsTiledGym04,
}

func (o *Onyx) GameplayScene() engine.SceneState {
	var cameraEntry *donburi.Entry
	var tilemapEntry *donburi.Entry
	var spriteEntry *donburi.Entry
	return engine.SceneState{
		OnEnter: func() error {
			assets := o.game.Assets()
			screen := o.game.Screen()
			tiledAssets := o.tiled.Assets()
			imageAssets := o.images.Assets()
			asepriteLibrary := o.aseprite.Library()

			screen.SetBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255})

			if err := assets.Load(content.AssetsFS(), gameplayManifest...); err != nil {
				return err
			}
			imgHandle, err := buildAnimations(assets, imageAssets, asepriteLibrary)
			if err != nil {
				return err
			}

			tilemap, tilemapHandle, err := buildTilemap(tiledAssets, content.AssetsTiledGym04)
			if err != nil {
				return err
			}
			tilemapCenter := tilemap.Bounds().Center()

			tilemapEntry = o.tiled.CreateTilemap(o.ecs.World(),
				tiled.WithTilemapHandle(tilemapHandle),
			)

			cameraEntry = transform.NewTransform(o.ecs.World())
			cameraEntry.AddComponent(camera.MainCamera)

			transform.SetPosition(cameraEntry, tilemapCenter.X, tilemapCenter.Y)
			camera.SetZoom(cameraEntry, 0.25)

			spriteEntry = o.aseprite.CreateSprite(o.ecs.World(),
				aseprite.WithImageOptions(
					images.WithHandle(imgHandle),
					images.WithAnchor(0.5, 1),
				),
				aseprite.WithClip("Idle"),
				aseprite.Playing(),
			)
			renderer.SetLayer(spriteEntry, 1)
			transform.SetPosition(spriteEntry, tilemapCenter.X, tilemapCenter.Y)

			bounds := transform.GetBounds(spriteEntry)
			width := bounds.Width() * 0.6
			height := bounds.Height() * 0.75

			collision.AddCollision(spriteEntry,
				collision.WithCollisionType(collision.CollisionTypeDynamic),
				collision.WithCollisionBounds(geom.AABB{
					Min: geom.Vec2{X: -width / 2, Y: -height},
					Max: geom.Vec2{X: width / 2, Y: 0},
				}),
			)

			o.Add(
				cameraEntry,
				tilemapEntry,
				spriteEntry,
			)
			return nil
		},
		OnUpdate: func(dt float64) (engine.SceneExitCode, error) {
			if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
				return engine.SceneExitNone, ebiten.Termination
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF) {
				ebiten.SetFullscreen(!ebiten.IsFullscreen())
			}

			move := geom.Vec2{}
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				move.Y -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				move.Y += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				move.X -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				move.X += 1
			}

			if move.X != 0 || move.Y != 0 {
				move = move.Normalize().Mul(100 * dt)
				transform.Translate(spriteEntry, move.X, move.Y)
				o.Update(spriteEntry)
			}

			return engine.SceneExitNone, nil
		},
		OnLateUpdate: func(dt float64) error {
			asepriteSystems := o.aseprite.Systems()

			viewport := camera.GetViewport(cameraEntry, o.game.Screen().SafeArea())
			o.ecs.QueryAll(viewport, func(entity donburi.Entity) {
				entry := o.ecs.World().Entry(entity)
				asepriteSystems.UpdateAnimation(entry, time.Duration(dt*float64(time.Second)))
			})
			return nil
		},
		OnExit: func() error {
			o.Remove(
				cameraEntry,
				tilemapEntry,
				spriteEntry,
			)
			return nil
		},
		OnRender: func(target *ebiten.Image) error {
			debug.DrawTransformBounds(o.ecs, cameraEntry, target, o.game.Screen().SafeArea())
			debug.DrawColliders(o.ecs.World(), o.collision.World(), cameraEntry, target, o.game.Screen().SafeArea())
			return nil
		},
	}
}

func buildAnimations(assets engine.Assets, imageAssets *images.ImageAssets, asepriteLibrary *aseprite.AsepriteLibrary) (uint64, error) {
	animationDataHandle, found := assets.GetDataHandle(content.AssetsAsepriteCaptainJson)
	if !found {
		return 0, engine.ErrAssetNotFound{Path: content.AssetsAsepriteCaptainJson.String()}
	}

	animationImageHandle, found := imageAssets.GetHandle(content.AssetsAsepriteCaptainImg)
	if !found {
		return 0, engine.ErrAssetNotFound{Path: content.AssetsAsepriteCaptainImg.String()}
	}

	animationData, _ := assets.GetData(animationDataHandle)
	asepriteLibrary.BuildAnimations(animationImageHandle, animationData)

	return animationImageHandle, nil
}

func buildTilemap(tiledAssets *tiled.TiledAssets, tmxPath file.FilePath) (*tiled.Tilemap, uint64, error) {
	tmxHandle, found := tiledAssets.GetTmxHandle(tmxPath)
	if !found {
		return nil, 0, engine.ErrAssetNotFound{Path: tmxPath.String()}
	}

	tilemap := tiledAssets.BuildTilemap(tmxHandle)
	return tilemap, tmxHandle, nil
}
