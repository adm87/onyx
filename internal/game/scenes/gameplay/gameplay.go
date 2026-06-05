package gameplay

import (
	"context"
	"fmt"
	"image/color"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
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

	debugDrawColliders := false
	debugDrawPartitions := false
	debugVisibilityToggle := true

	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			collision.AddCollisionEnter(world, onCollisionEnter)
			collision.AddCollisionExit(world, onCollisionExit)

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
			colliders.AddCollider(entry,
				colliders.WithAABB(
					geom.AABB{
						Min: geom.Vec2{X: -hWidth, Y: -hHeight},
						Max: geom.Vec2{X: hWidth, Y: hHeight},
					},
				),
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

			collision.Add(entry)
			player = entry.Entity()

			return nil
		},
		OnExit: func(ctx context.Context, world donburi.World) error {
			collision.RemoveCollisionEnter(world, onCollisionEnter)
			collision.RemoveCollisionExit(world, onCollisionExit)
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

			transform.SetPosition(entry, position)
			collision.Update(entry)

			camera.SetPosition(world, position)

			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world donburi.World, img *ebiten.Image, viewMatrix ebiten.GeoM) error {
			colliders.Query(world, func(entry *donburi.Entry, layer colliders.CollisionLayer, aabb geom.AABB) {
				aabb = aabb.Translate(transform.GetPosition(entry))
				center := camera.ToScreen(world, screen, aabb.Center())

				if screen.SafeArea().Contains(center) {
					min := camera.ToScreen(world, screen, aabb.Min)
					max := camera.ToScreen(world, screen, aabb.Max)
					vector.FillRect(
						img,
						float32(min.X),
						float32(min.Y),
						float32(max.X-min.X),
						float32(max.Y-min.Y),
						color.RGBA{G: 255, A: 255},
						false,
					)
					vector.FillRect(
						img,
						float32(center.X)-2,
						float32(center.Y)-2,
						4,
						4,
						color.RGBA{R: 255, A: 255},
						false,
					)
					ebitenutil.DebugPrintAt(img, fmt.Sprintf("%d", entry.Entity()), int(center.X), int(center.Y))
				}
			})
			return nil
		},
	}
}

func buildStaticCollision(world donburi.World, collision engine.Collision, tmx *tiled.Tmx) {
	tmx.ObjectGroups.EachInGroup("collision_static", func(object *tiled.TmxObject) {
		entry := colliders.NewCollider(world,
			colliders.AsStatic(),
			colliders.WithAABB(
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

func onCollisionEnter(world donburi.World, event engine.CollisionEvent) {
	fmt.Printf("Collision Enter: %d <-> %d\n", event.EntityA, event.EntityB)
}

func onCollisionExit(world donburi.World, event engine.CollisionEvent) {
	fmt.Printf("Collision Exit: %d <-> %d\n", event.EntityA, event.EntityB)
}
