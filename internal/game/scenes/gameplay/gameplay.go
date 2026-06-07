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
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
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

	var moveDir geom.Vec2
	var zoomDir int

	var manifest = []file.FilePath{
		content.AssetsAsepriteCaptainJson,
		content.AssetsAsepriteCaptainImg,
		content.AssetsTiledGym01,
	}
	var animations = []file.FilePath{
		content.AssetsAsepriteCaptainJson,
	}

	debugDrawColliders := false
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

			playerEntry, err = buildPlayer(ecs, assets, tilemap.Bounds().Center())
			if err != nil {
				return fmt.Errorf("failed to build player: %w", err)
			}

			world.Add(playerEntry)
			world.Add(tilemapEntry)

			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				moveDir.Y = -1
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				moveDir.Y = 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				moveDir.X = -1
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
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

			return engine.SceneExitNone, nil
		},
		OnFixedUpdate: func(ecs donburi.World, dt float64) error {
			position := transform.GetPosition(playerEntry)
			if moveDir.X != 0 || moveDir.Y != 0 {

				position.X += moveDir.X * 100 * dt
				position.Y += moveDir.Y * 100 * dt

				viewport := camera.Viewport()

				position.X = engine.Clamp(position.X, viewport.Min.X, viewport.Max.X)
				position.Y = engine.Clamp(position.Y, viewport.Min.Y, viewport.Max.Y)

				transform.SetPosition(playerEntry, position)

				moveDir = geom.Vec2{}
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

			cameraPosition.X = engine.SmoothStep(cameraPosition.X, playerPosition.X, dt*10)
			cameraPosition.Y = engine.SmoothStep(cameraPosition.Y, playerPosition.Y, dt*10)

			if math.Abs(cameraPosition.X-playerPosition.X) < 0.1 {
				cameraPosition.X = playerPosition.X
			}
			if math.Abs(cameraPosition.Y-playerPosition.Y) < 0.1 {
				cameraPosition.Y = playerPosition.Y
			}

			cameraPosition.X = engine.Clamp(cameraPosition.X, min.X+(width/2), max.X-(width/2))
			cameraPosition.Y = engine.Clamp(cameraPosition.Y, min.Y+(height/2), max.Y-(height/2))

			camera.SetPosition(cameraPosition)

			world.Update(playerEntry)
			return nil
		},
		OnRender: func(ecs donburi.World, img *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error {
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

func buildPlayer(ecs donburi.World, assets engine.Assets, spawnPosition geom.Vec2) (*donburi.Entry, error) {
	img, exists := images.GetImageAssets(assets, content.EmbeddedImg10x10White)
	if !exists {
		return nil, fmt.Errorf("failed to load embedded image: %s", content.EmbeddedImg10x10White)
	}

	width, height := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	hWidth, hHeight := width/2, height/2

	playerBounds := geom.AABB{
		Min: geom.Vec2{X: -hWidth, Y: -hHeight},
		Max: geom.Vec2{X: hWidth, Y: hHeight},
	}

	entry := images.CreateImageEntity(ecs, content.EmbeddedImg10x10White, playerBounds)

	colliders.AddCollider(entry)

	rendering.SetAnchor(entry, geom.Vec2{X: 0.5, Y: 0.5})
	rendering.SetLayer(entry, 1)

	transform.SetPosition(entry, spawnPosition)

	return entry, nil
}

func buildTiledLevel(ecs donburi.World, assets engine.Assets, world engine.World) (*donburi.Entry, *tiled.Tilemap, error) {
	tmx, exists := tiled.GetTmx(assets, content.AssetsTiledGym01)
	if !exists {
		return nil, nil, fmt.Errorf("tmx asset not found for tilemap: %s", content.AssetsTiledGym01)
	}

	tilemap, exists := tiled.GetTilemap(assets, content.AssetsTiledGym01)
	if !exists {
		return nil, nil, fmt.Errorf("tilemap asset not found: %s", content.AssetsTiledGym01)
	}

	levelEntry := tiled.CreateTiledEntity(ecs, content.AssetsTiledGym01, tilemap.Bounds())
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
