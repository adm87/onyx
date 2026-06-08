package gameplay

import (
	"fmt"
	"math"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/aseprite"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/shapes"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

const (
	PlayerIdleAnim = "Idle"
	PlayerRunAnim  = "Run"
)

func New(
	assets engine.Assets,
	camera engine.Camera,
	collision engine.Collision,
	screen engine.Screen,
	world engine.World) engine.SceneState {

	var asepriteAdapter *aseprite.AsepriteAssetAdapter
	var tilemap *tiled.Tilemap
	var tilemapEntry *donburi.Entry
	var playerEntry *donburi.Entry

	var path vector.Path
	var moveDir geom.Vec2
	var zoomDir int

	var manifest = []file.FilePath{
		content.AssetsAsepriteCaptain,
		content.AssetsAsepriteCaptainImg,
		content.AssetsTiledGym04,
	}
	var animations = []file.FilePath{
		content.AssetsAsepriteCaptain,
	}

	debugDrawColliders := true
	debugDrawPartitions := false
	debugVisibilityToggle := true

	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			var err error
			var found bool

			if err = assets.Load(content.AssetsFS(), manifest...); err != nil {
				return fmt.Errorf("failed to load assets: %v", err)
			}

			asepriteAdapter, found = aseprite.GetAssetAdapter(assets)
			if !found {
				return fmt.Errorf("aseprite asset adapter not found")
			}

			if err = asepriteAdapter.ImportAnimations(animations...); err != nil {
				return fmt.Errorf("failed to import aseprite animations: %w", err)
			}

			tilemapEntry, tilemap, err = buildTiledLevel(ecs, assets, world)
			if err != nil {
				return fmt.Errorf("failed to build tiled level: %w", err)
			}

			viewport := camera.Viewport()

			scaleX := tilemap.Bounds().Width() / viewport.Width()
			scaleY := tilemap.Bounds().Height() / viewport.Height()
			zoom := min(scaleX, scaleY)

			camera.SetPosition(tilemap.Bounds().Center())
			camera.SetZoom(zoom)

			playerEntry, err = buildPlayer(ecs, tilemap.Bounds().Center())
			if err != nil {
				return fmt.Errorf("failed to build player: %w", err)
			}

			world.Add(playerEntry)
			world.Add(tilemapEntry)

			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				moveDir.X = -1
			} else if ebiten.IsKeyPressed(ebiten.KeyD) {
				moveDir.X = 1
			}

			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				zoomDir = 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) {
				zoomDir = -1
			}

			if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
				return engine.SceneExitNone, ebiten.Termination
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF) {
				ebiten.SetFullscreen(!ebiten.IsFullscreen())
			}

			if inpututil.IsKeyJustPressed(ebiten.Key1) {
				debugDrawColliders = !debugDrawColliders
			}
			if inpututil.IsKeyJustPressed(ebiten.Key2) {
				debugDrawPartitions = !debugDrawPartitions
			}
			if inpututil.IsKeyJustPressed(ebiten.Key3) {
				debugVisibilityToggle = !debugVisibilityToggle
				rendering.SetVisible(tilemapEntry, debugVisibilityToggle)
				rendering.SetVisible(playerEntry, debugVisibilityToggle)
			}

			if moveDir.X != 0 {
				transform.SetScale(playerEntry, geom.Vec2{X: moveDir.X, Y: 1})
				aseprite.SetAnimationName(playerEntry, PlayerRunAnim)
			} else {
				aseprite.SetAnimationName(playerEntry, PlayerIdleAnim)
			}

			updateRegion := camera.Viewport().Scale(2)
			world.Query(updateRegion, func(e donburi.Entity) {
				entry := ecs.Entry(e)

				if entry.HasComponent(aseprite.Animation) {
					aseprite.UpdateAnimation(entry, asepriteAdapter, dt)
				}
			})

			return engine.SceneExitNone, nil
		},
		OnFixedUpdate: func(ecs donburi.World, dt float64) error {
			position := transform.GetPosition(playerEntry)
			if moveDir.X != 0 || moveDir.Y != 0 {

				position.X += moveDir.X * 100 * dt
				position.Y += moveDir.Y * 100 * dt

				viewport := camera.Viewport()
				aabb := shapes.GetAABB(playerEntry).Translate(position)

				if aabb.Min.X < viewport.Min.X {
					position.X = viewport.Min.X - aabb.Min.X + position.X
				}
				if aabb.Max.X > viewport.Max.X {
					position.X = viewport.Max.X - aabb.Max.X + position.X
				}
				if aabb.Min.Y < viewport.Min.Y {
					position.Y = viewport.Min.Y - aabb.Min.Y + position.Y
				}
				if aabb.Max.Y > viewport.Max.Y {
					position.Y = viewport.Max.Y - aabb.Max.Y + position.Y
				}

				moveDir = geom.Vec2{X: 0, Y: 0}

				transform.SetPosition(playerEntry, position)
			}
			if zoomDir != 0 {
				applyAndClampZoom(camera, tilemap.Bounds(), float64(zoomDir)*0.1*dt)
				zoomDir = 0
			}
			return nil
		},
		OnLateUpdate: func(ecs donburi.World, dt float64) error {
			playerPosition := transform.GetPosition(playerEntry)
			cameraPosition := camera.Position()

			min := tilemap.Bounds().Min
			max := tilemap.Bounds().Max

			viewport := camera.Viewport()
			width := viewport.Width()
			height := viewport.Height()

			cameraPosition.X = engine.SmoothStep(cameraPosition.X, playerPosition.X, dt*15)
			cameraPosition.Y = engine.SmoothStep(cameraPosition.Y, playerPosition.Y, dt*15)

			if math.Abs(cameraPosition.X-playerPosition.X) < 0.05 {
				cameraPosition.X = playerPosition.X
			}
			if math.Abs(cameraPosition.Y-playerPosition.Y) < 0.05 {
				cameraPosition.Y = playerPosition.Y
			}

			cameraPosition.X = engine.Clamp(cameraPosition.X, min.X+(width/2), max.X-(width/2))
			cameraPosition.Y = engine.Clamp(cameraPosition.Y, min.Y+(height/2), max.Y-(height/2))

			camera.SetPosition(cameraPosition)

			world.Update(playerEntry)
			return nil
		},
		OnRender: func(ecs donburi.World, img *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error {
			path.Reset()

			if debugDrawColliders {
				world.Query(viewport, func(e donburi.Entity) {
					entry := ecs.Entry(e)

					aabb := shapes.GetAABB(entry).Translate(transform.GetPosition(entry))

					screenMin := camera.ToScreen(aabb.Min)
					screenMax := camera.ToScreen(aabb.Max)

					path.MoveTo(float32(screenMin.X), float32(screenMin.Y))
					path.LineTo(float32(screenMax.X), float32(screenMin.Y))
					path.LineTo(float32(screenMax.X), float32(screenMax.Y))
					path.LineTo(float32(screenMin.X), float32(screenMax.Y))
					path.Close()
				})

				vector.StrokePath(img, &path, &vector.StrokeOptions{
					Width: 2,
				}, &vector.DrawPathOptions{})
			}

			ebitenutil.DebugPrintAt(img, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 10, 10)
			return nil
		},
	}
}

func applyAndClampZoom(camera engine.Camera, bounds geom.AABB, appliedZoom float64) {
	zoom := engine.Clamp(camera.Zoom()+appliedZoom, 0.1, 10.0)
	camera.SetZoom(zoom)

	viewport := camera.Viewport()
	if viewport.Width() > bounds.Width() {
		camera.SetZoom(camera.Zoom() - appliedZoom)
	}
	if viewport.Height() > bounds.Height() {
		camera.SetZoom(camera.Zoom() - appliedZoom)
	}
}

func buildPlayer(ecs donburi.World, spawnPosition geom.Vec2) (*donburi.Entry, error) {
	entry := aseprite.CreateSprite(
		ecs,
		content.AssetsAsepriteCaptain,
		PlayerIdleAnim,
		geom.AABB{
			Min: geom.Vec2{X: -5, Y: -15},
			Max: geom.Vec2{X: 5, Y: 0},
		},
	)

	colliders.AddCollider(entry)

	rendering.SetLayer(entry, 1)
	rendering.SetAnchor(entry, geom.Vec2{X: 0.5, Y: 1.0})

	transform.SetPosition(entry, spawnPosition)

	return entry, nil
}

func buildTiledLevel(ecs donburi.World, assets engine.Assets, world engine.World) (*donburi.Entry, *tiled.Tilemap, error) {
	tmx, exists := tiled.GetTmx(assets, content.AssetsTiledGym04)
	if !exists {
		return nil, nil, fmt.Errorf("tmx asset not found for tilemap: %s", content.AssetsTiledGym04)
	}

	tilemap, exists := tiled.GetTilemap(assets, content.AssetsTiledGym04)
	if !exists {
		return nil, nil, fmt.Errorf("tilemap asset not found: %s", content.AssetsTiledGym04)
	}

	levelEntry := tiled.CreateTiledMap(ecs, content.AssetsTiledGym04, tilemap.Bounds())
	rendering.SetLayer(levelEntry, 0)

	buildStaticLevelCollision(ecs, world, tmx)
	return levelEntry, tilemap, nil
}

func buildStaticLevelCollision(ecs donburi.World, world engine.World, tmx *tiled.Tmx) {
	tmx.ObjectGroups.EachInGroup("collision_static", func(object *tiled.TmxObject) {
		position := geom.Vec2{
			X: object.X,
			Y: object.Y,
		}
		entry := colliders.NewCollider(ecs,
			colliders.AsStatic(),
		)
		shapes.AddAABB(entry,
			shapes.WithPosition(position.X, position.Y),
			shapes.WithSize(object.Width, object.Height),
		)
		world.Add(entry)
	})
}
