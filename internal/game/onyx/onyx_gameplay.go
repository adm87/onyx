package onyx

import (
	"image/color"
	"time"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/internal/game/components/movement"
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

var (
	debugDrawTransformBounds  = false
	debugDrawColliders        = false
	debugDrawNearestColliders = false
)

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

			tilemap, tilemapHandle, err := buildTilemap(o, tiledAssets, content.AssetsTiledGym04)
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
				collision.WithCollisionBounds(
					geom.Vec2{X: -width / 2, Y: -height},
					geom.Vec2{X: width / 2, Y: 0},
				),
			)

			movement.AddMovement(spriteEntry,
				movement.WithSpeed(100),
			)

			o.AddEntries(
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
			if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
				debugDrawTransformBounds = !debugDrawTransformBounds
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
				debugDrawColliders = !debugDrawColliders
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
				debugDrawNearestColliders = !debugDrawNearestColliders
			}

			var moveX, moveY float64

			movement.ClearDirection(spriteEntry)
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				moveY -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				moveY += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				moveX -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				moveX += 1
			}
			movement.SetDirection(spriteEntry, moveX, moveY)

			return engine.SceneExitNone, nil
		},
		OnFixedUpdate: func(dt float64) error {
			movement.ApplyMovement(o.ecs.World(), dt)
			if movement.IsMoving(spriteEntry) {
				direction := movement.GetDirection(spriteEntry)
				if direction.X < 0 {
					transform.SetScale(spriteEntry, -1, 1)
				} else if direction.X > 0 {
					transform.SetScale(spriteEntry, 1, 1)
				}
				o.UpdateEntries(spriteEntry)
			}

			collisionSystems := o.collision.Systems()

			spriteCollider := collision.GetWorldCollider(spriteEntry)
			if infos, ok := collisionSystems.CheckStaticCollision(o.ecs.World(), spriteCollider); ok {
				_ = infos
			}
			return nil
		},
		OnLateUpdate: func(dt float64) error {
			direction := movement.GetDirection(spriteEntry)
			if direction.X == 0 {
				aseprite.SetClip(spriteEntry, "Idle")
			} else {
				aseprite.SetClip(spriteEntry, "Run")
			}

			asepriteSystems := o.aseprite.Systems()

			viewport := camera.GetViewport(cameraEntry, o.game.Screen().SafeArea())
			o.ecs.QueryAll(viewport, func(entity donburi.Entity) {
				entry := o.ecs.World().Entry(entity)
				asepriteSystems.UpdateAnimation(entry, time.Duration(dt*float64(time.Second)))
			})
			return nil
		},
		OnExit: func() error {
			o.RemoveEntries(
				cameraEntry,
				tilemapEntry,
				spriteEntry,
			)
			return nil
		},
		OnRender: func(target *ebiten.Image) error {
			safeArea := o.game.Screen().SafeArea()
			if debugDrawTransformBounds {
				debug.DrawTransformBounds(o.ecs, cameraEntry, target, safeArea)
			}
			if debugDrawColliders {
				debug.DrawColliders(
					o.ecs.World(),
					camera.GetViewport(cameraEntry, safeArea),
					o.collision.World(),
					cameraEntry,
					target,
					safeArea,
					color.RGBA{R: 255, A: 255},
				)
			}
			if debugDrawNearestColliders {
				aabb := collision.GetWorldCollider(spriteEntry)
				debug.DrawStaticPartitioner(
					o.collision.World(),
					cameraEntry,
					target,
					aabb,
					safeArea,
				)
				debug.DrawColliders(
					o.ecs.World(),
					aabb,
					o.collision.World(),
					cameraEntry,
					target,
					safeArea,
					color.RGBA{G: 255, A: 255},
				)
			}
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

func buildTilemap(game *Onyx, tiledAssets *tiled.TiledAssets, tmxPath file.FilePath) (*tiled.Tilemap, uint64, error) {
	tmxHandle, found := tiledAssets.GetTmxHandle(tmxPath)
	if !found {
		return nil, 0, engine.ErrAssetNotFound{Path: tmxPath.String()}
	}

	tilemap, tmx := tiledAssets.BuildTilemap(tmxHandle)
	tmx.ObjectGroups.EachInGroup("collision", func(object *tiled.TmxObject) {
		min := geom.Vec2{}
		max := geom.Vec2{X: object.Width, Y: object.Height}

		entry := transform.NewTransform(game.ecs.World(),
			transform.WithPosition(object.X, object.Y),
			transform.WithBounds(min, max),
		)

		collision.AddCollision(entry,
			collision.WithCollisionBounds(min, max),
		)

		game.AddEntries(entry)
	})

	return tilemap, tmxHandle, nil
}
