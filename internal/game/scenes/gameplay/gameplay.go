package gameplay

import (
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/shapes"
	"github.com/adm87/onyx/pkg/engine/components/transform"
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

	const tilemapRef = content.AssetsLevelsGym01

	var player donburi.Entity
	var level donburi.Entity
	var tilemap *tiled.Tilemap
	var moveDir geom.Vec2
	var zoomDir int

	debugDrawColliders := false
	debugDrawPartitions := false
	debugVisibilityToggle := true

	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			if err := assets.Load(content.AssetsFS(), tilemapRef); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}

			tm, exists := tiled.GetTilemap(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tilemap asset not found: %s", tilemapRef)
			}

			tilemap = tm

			tmx, exists := tiled.GetTmx(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tmx asset not found for tilemap: %s", tilemapRef)
			}

			img, exists := images.GetImageAssets(assets, content.EmbeddedImg10x10White)
			if !exists {
				return fmt.Errorf("failed to load embedded image: %s", content.EmbeddedImg10x10White)
			}

			levelEntry := tiled.CreateTiledEntity(ecs, content.AssetsLevelsGym01, tilemap.Bounds())
			level = levelEntry.Entity()

			rendering.SetLayer(levelEntry, 0)

			buildStaticCollision(ecs, world, tmx)

			camera.SetPosition(tilemap.Bounds().Center())
			camera.SetZoom(0.2)

			width, height := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
			hWidth, hHeight := width/2, height/2

			playerBounds := geom.AABB{
				Min: geom.Vec2{X: -hWidth, Y: -hHeight},
				Max: geom.Vec2{X: hWidth, Y: hHeight},
			}

			playerEntry := images.CreateImageEntity(ecs, content.EmbeddedImg10x10White, playerBounds)

			colliders.AddCollider(playerEntry)
			rendering.SetAnchor(playerEntry, geom.Vec2{X: 0.5, Y: 0.5})
			rendering.SetLayer(playerEntry, 1)

			pos := tilemap.Bounds().Center()
			transform.SetPosition(playerEntry, pos)

			player = playerEntry.Entity()

			world.Add(playerEntry)
			world.Add(levelEntry)

			return nil
		},
		OnExit: func(ecs donburi.World) error {
			world.Remove(ecs.Entry(player))
			world.Remove(ecs.Entry(level))
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
				rendering.SetVisible(ecs.Entry(level), debugVisibilityToggle)
				rendering.SetVisible(ecs.Entry(player), debugVisibilityToggle)
			}

			return engine.SceneExitNone, nil
		},
		OnFixedUpdate: func(ecs donburi.World, dt float64) error {
			if moveDir.X != 0 || moveDir.Y != 0 {
				entry := ecs.Entry(player)

				position := transform.GetPosition(entry)

				position.X += moveDir.X * 100 * dt
				position.Y += moveDir.Y * 100 * dt

				viewport := camera.Viewport()

				position.X = engine.Clamp(position.X, viewport.Min.X, viewport.Max.X)
				position.Y = engine.Clamp(position.Y, viewport.Min.Y, viewport.Max.Y)

				transform.SetPosition(entry, position)

				moveDir = geom.Vec2{}
			}
			if zoomDir != 0 {
				zoom := camera.Zoom()
				zoom *= 1 + (0.5 * float64(zoomDir) * dt)
				camera.SetZoom(zoom)
				zoomDir = 0
			}
			return nil
		},
		OnLateUpdate: func(ecs donburi.World, dt float64) error {
			playerEntry := ecs.Entry(player)

			position := transform.GetPosition(playerEntry)

			min := tilemap.Bounds().Min
			max := tilemap.Bounds().Max

			viewport := camera.Viewport()
			width := viewport.Width()
			height := viewport.Height()

			position.X = engine.Clamp(position.X, min.X+(width/2), max.X-(width/2))
			position.Y = engine.Clamp(position.Y, min.Y+(height/2), max.Y-(height/2))

			camera.SetPosition(position)

			world.Update(playerEntry)
			return nil
		},
		OnRender: func(ecs donburi.World, img *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error {

			ebitenutil.DebugPrintAt(img, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 10, 10)
			return nil
		},
	}
}

func buildStaticCollision(ecs donburi.World, world engine.World, tmx *tiled.Tmx) {
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
