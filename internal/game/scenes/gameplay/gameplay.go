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
	screen engine.Screen,
	time engine.Time) engine.SceneState {

	const tilemapRef = content.AssetsLevelsGym01

	var player donburi.Entity
	var level donburi.Entity

	debugDrawColliders := false
	debugDrawPartitions := false
	debugVisibilityToggle := true

	return engine.SceneState{
		OnEnter: func(ctx context.Context, world engine.World) error {
			collision := world.Collision()
			ecs := world.ECS()

			collision.AddCollisionEnter(ecs, onCollisionEnter)
			collision.AddCollisionExit(ecs, onCollisionExit)

			if err := assets.Load(content.AssetsFS(), tilemapRef); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}

			lvlEntry := tiled.CreateTiledEntity(ecs, content.AssetsLevelsGym01)
			level = lvlEntry.Entity()

			tilemap, exists := tiled.GetTilemap(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tilemap asset not found: %s", tilemapRef)
			}

			tmx, exists := tiled.GetTmx(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tmx asset not found for tilemap: %s", tilemapRef)
			}

			buildStaticCollision(world, tmx)

			camera.SetPosition(tilemap.Bounds().Center())
			camera.SetZoom(0.2)

			img, exists := images.GetImageAssets(assets, content.EmbeddedImg10x10White)
			if !exists {
				return fmt.Errorf("failed to load embedded image: %s", content.EmbeddedImg10x10White)
			}

			width, height := img.Bounds().Dx(), img.Bounds().Dy()
			hWidth, hHeight := float64(width)/2, float64(height)/2

			aabb := geom.AABB{
				Min: geom.Vec2{X: -hWidth, Y: -hHeight},
				Max: geom.Vec2{X: hWidth, Y: hHeight},
			}

			entry := images.CreateImageEntity(ecs, content.EmbeddedImg10x10White)
			colliders.AddCollider(entry,
				colliders.WithAABB(aabb),
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

			player = entry.Entity()

			world.Add(entry)
			return nil
		},
		OnExit: func(ctx context.Context, world engine.World) error {
			collision := world.Collision()
			ecs := world.ECS()

			collision.RemoveCollisionEnter(ecs, onCollisionEnter)
			collision.RemoveCollisionExit(ecs, onCollisionExit)
			return nil
		},
		OnUpdate: func(ctx context.Context, world engine.World) (engine.SceneExitCode, error) {
			ecs := world.ECS()

			entry := ecs.Entry(player)
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
				zoom := camera.Zoom()
				zoom *= 1 + (0.5 * time.DeltaTime())
				camera.SetZoom(zoom)
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) {
				zoom := camera.Zoom()
				zoom /= 1 + (0.5 * time.DeltaTime())
				camera.SetZoom(zoom)
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

			transform.SetPosition(entry, position)
			camera.SetPosition(position)

			world.Update(entry)
			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world engine.World, img *ebiten.Image, viewMatrix ebiten.GeoM) error {
			ecs := world.ECS()
			collision := world.Collision()

			min := camera.ToWorld(screen, screen.SafeArea().Min)
			max := camera.ToWorld(screen, screen.SafeArea().Max)
			queryRegion := geom.AABB{Min: min, Max: max}

			count := 0
			collision.QueryAll(queryRegion, func(entity donburi.Entity) {
				count++
				entry := ecs.Entry(entity)

				position := transform.GetPosition(entry)
				aabb := colliders.GetAABB(entry).Translate(position)

				center := camera.ToScreen(screen, aabb.Center())

				vector.FillRect(
					img,
					float32(center.X)-2,
					float32(center.Y)-2,
					4,
					4,
					color.RGBA{R: 255, A: 255},
					false,
				)

				ebitenutil.DebugPrintAt(img, fmt.Sprintf("%d", entity), int(center.X), int(center.Y))
			})

			ebitenutil.DebugPrintAt(img, fmt.Sprintf("Colliders in view: %d", count), 10, 10)
			return nil
		},
	}
}

func buildStaticCollision(world engine.World, tmx *tiled.Tmx) {
	ecs := world.ECS()

	tmx.ObjectGroups.EachInGroup("collision_static", func(object *tiled.TmxObject) {
		aabb := geom.AABB{
			Max: geom.Vec2{
				X: object.Width,
				Y: object.Height,
			},
		}
		position := geom.Vec2{
			X: object.X,
			Y: object.Y,
		}
		entry := colliders.NewCollider(ecs,
			colliders.AsStatic(),
			colliders.WithAABB(aabb),
		)
		transform.SetPosition(entry, position)
		world.Add(entry)
	})
}

func onCollisionEnter(ecs donburi.World, event engine.CollisionEvent) {
	fmt.Printf("Collision Enter: %d <-> %d\n", event.EntityA, event.EntityB)
}

func onCollisionExit(ecs donburi.World, event engine.CollisionEvent) {
	fmt.Printf("Collision Exit: %d <-> %d\n", event.EntityA, event.EntityB)
}
