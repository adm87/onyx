package gameplay

import (
	"context"
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

const (
	defaultCollisionLayer colliders.CollisionLayer = 1 << iota
)

func New(
	assets engine.Assets,
	camera engine.Camera,
	collision engine.Collision,
	screen engine.Screen,
	time engine.Time) engine.SceneState {

	const tilemapRef = content.AssetsLevelsGym01

	var player donburi.Entity
	var level donburi.Entity
	var corners [4]donburi.Entity

	debugDrawColliders := false
	debugDrawPartitions := false
	debugVisibilityToggle := true

	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			// engine.StaticCollision.OnEnter.Subscribe(world, onCollisionEnter)
			// engine.StaticCollision.OnExit.Subscribe(world, onCollisionExit)

			if err := assets.Load(content.AssetsFS(), tilemapRef); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}

			lvlEntry := tiled.CreateTiledEntity(world, content.AssetsLevelsGym01)
			level = lvlEntry.Entity()

			tilemap, exists := tiled.GetTilemap(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tilemap asset not found: %s", tilemapRef)
			}

			tmx, exists := tiled.GetTmx(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tmx asset not found for tilemap: %s", tilemapRef)
			}

			buildStaticCollision(world, collision, tmx)

			camera.SetPosition(world, tilemap.Bounds().Center())
			camera.SetZoom(world, 0.2)

			img, exists := images.GetImageAssets(assets, content.EmbeddedImg10x10White)
			if !exists {
				return fmt.Errorf("failed to load embedded image: %s", content.EmbeddedImg10x10White)
			}

			width, height := img.Bounds().Dx(), img.Bounds().Dy()
			hWidth, hHeight := float64(width)/2, float64(height)/2

			entry := images.CreateImageEntity(world, content.EmbeddedImg10x10White)
			colliders.AddBoxCollider(entry,
				colliders.WithBox(
					geom.AABB{
						Min: geom.Vec2{X: -hWidth, Y: -hHeight},
						Max: geom.Vec2{X: hWidth, Y: hHeight},
					},
				),
				colliders.WithLayer(defaultCollisionLayer),
			)
			rendering.SetAnchor(entry,
				geom.Vec2{
					X: 0.5,
					Y: 0.5,
				},
			)
			rendering.SetLayer(entry, 1)

			pos := tilemap.Bounds().Center()
			transform.SetPosition(entry, pos)

			for i := range corners {
				aX, aY := 0.0, 0.0
				switch i {
				case 0:
					aX, aY = 0, 0
				case 1:
					aX, aY = 1, 0
				case 2:
					aX, aY = 1, 1
				case 3:
					aX, aY = 0, 1
				}

				cornerEntry := images.CreateImageEntity(world, content.EmbeddedImg10x10White)
				rendering.SetZIndex(cornerEntry, 1)
				rendering.SetAnchor(cornerEntry, geom.Vec2{X: aX, Y: aY})

				corners[i] = cornerEntry.Entity()
			}

			collision.Add(entry)
			player = entry.Entity()

			return nil
		},
		OnExit: func(ctx context.Context, world donburi.World) error {
			// engine.StaticCollision.OnEnter.Unsubscribe(world, onCollisionEnter)
			// engine.StaticCollision.OnExit.Unsubscribe(world, onCollisionExit)
			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			entry := world.Entry(player)
			position := transform.GetPosition(entry)

			if ebiten.IsKeyPressed(ebiten.KeyW) {
				position.Y -= 100 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				position.Y += 100 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				position.X -= 100 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				position.X += 100 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				zoom := camera.Zoom(world)
				zoom *= 1 + (0.5 * time.DeltaTime())
				camera.SetZoom(world, zoom)
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) {
				zoom := camera.Zoom(world)
				zoom /= 1 + (0.5 * time.DeltaTime())
				camera.SetZoom(world, zoom)
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
				rendering.SetVisible(world.Entry(level), debugVisibilityToggle)
				rendering.SetVisible(world.Entry(player), debugVisibilityToggle)
			}

			camera.SetPosition(world, position)
			transform.SetPosition(entry, position)

			min := camera.ToWorld(world, screen, screen.SafeArea().Min)
			max := camera.ToWorld(world, screen, screen.SafeArea().Max)

			for i := range corners {
				entry := world.Entry(corners[i])
				switch i {
				case 0:
					transform.SetPosition(entry, min)
				case 1:
					transform.SetPosition(entry, geom.Vec2{X: max.X, Y: min.Y})
				case 2:
					transform.SetPosition(entry, max)
				case 3:
					transform.SetPosition(entry, geom.Vec2{X: min.X, Y: max.Y})
				}
			}

			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world donburi.World, img *ebiten.Image, viewMatrix ebiten.GeoM) error {
			return nil
		},
	}
}

func buildStaticCollision(world donburi.World, collision engine.Collision, tmx *tiled.Tmx) {
	tmx.ObjectGroups.EachInGroup("collision_static", func(object *tiled.TmxObject) {
		entry := colliders.NewBoxCollider(world,
			colliders.AsStatic(),
			colliders.WithBox(
				geom.AABB{
					Min: geom.Vec2{X: 0, Y: 0},
					Max: geom.Vec2{X: object.Width, Y: object.Height},
				},
			),
		)
		transform.SetPosition(entry, geom.Vec2{X: object.X, Y: object.Y})
		collision.Add(entry)
	})
}

// func onCollisionEnter(world donburi.World, event engine.CollisionEvent) {
// 	now := time.Now()
// 	println("Collision Enter:", event.EntityA, "with", event.EntityB, "at", now.Format("15:04:05.000"))
// }

// func onCollisionExit(world donburi.World, event engine.CollisionEvent) {
// 	now := time.Now()
// 	println("Collision Exit:", event.EntityA, "with", event.EntityB, "at", now.Format("15:04:05.000"))
// }
