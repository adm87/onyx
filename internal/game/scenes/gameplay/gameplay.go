package gameplay

import (
	"context"
	"fmt"
	"image/color"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/colliders"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/engine/partitioning"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
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

	var entity donburi.Entity
	var path vector.Path

	debugDrawColliders := false
	debugDrawPartitions := false

	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			if err := assets.Load(content.AssetsFS(), tilemapRef); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}
			tiled.CreateTilemap(world,
				tiled.WithTilemapRef(tilemapRef),
			)

			tilemap, exists := tiled.GetTilemap(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tilemap asset not found: %s", tilemapRef)
			}

			tmx, exists := tiled.GetTmx(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tmx asset not found for tilemap: %s", tilemapRef)
			}

			buildStaticCollision(world, collision, tmx)

			camera.SetPosition(tilemap.Bounds().Center())
			camera.SetZoom(0.2)

			img, exists := images.GetImage(assets, content.EmbeddedImg10x10White)
			if !exists {
				return fmt.Errorf("failed to load embedded image: %s", content.EmbeddedImg10x10White)
			}

			width, height := img.Bounds().Dx(), img.Bounds().Dy()
			halfWidth, halfHeight := float64(width)/2, float64(height)/2

			entry := images.CreateImageEntity(world,
				images.WithRef(content.EmbeddedImg10x10White),
				images.WithLayer(1),
				images.WithPosition(tilemap.Bounds().Center().XY()),
				images.WithAnchor(0.5, 0.5),
			)

			colliders.SetBoxCollider(entry, geom.AABB{
				Min: geom.Vec2{X: -halfWidth, Y: -halfHeight},
				Max: geom.Vec2{X: halfWidth, Y: halfHeight},
			})
			colliders.SetColliderType(entry, colliders.ColliderTypeDynamic)

			collision.Add(entry)

			entity = entry.Entity()
			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			entry := world.Entry(entity)
			position := transform.GetPosition(entry)

			if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
				position.Y -= 100 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
				position.Y += 100 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
				position.X -= 100 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
				position.X += 100 * time.DeltaTime()
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

			camera.SetPosition(position)
			transform.SetPosition(entry, position)

			collision.Update(entry)

			collision.Simulate()
			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world donburi.World, img *ebiten.Image, viewMatrix ebiten.GeoM) error {
			if debugDrawColliders {
				debugDrawEntityColliders(world, camera, screen, collision, &path, img, viewMatrix)
			}
			if debugDrawPartitions {
				partitioning.DebugDrawSpatialHash(img, collision.Partitioning(), screen.SafeArea(), viewMatrix, color.RGBA{G: 255, A: 255})
				partitioning.DebugDrawSpatialHash(img, collision.StaticPartitioning(), screen.SafeArea(), viewMatrix, color.RGBA{R: 255, A: 255})
			}
			return nil
		},
	}
}

func buildStaticCollision(world donburi.World, collision engine.Collision, tmx *tiled.Tmx) {
	tmx.ObjectGroups.EachInGroup("collision_static", func(object *tiled.TmxObject) {
		entry := world.Entry(world.Create(
			colliders.BoxCollider,
			transform.Position,
		))
		colliders.SetBoxCollider(entry, geom.AABB{
			Min: geom.Vec2{X: 0, Y: 0},
			Max: geom.Vec2{X: object.Width, Y: object.Height},
		})
		colliders.SetColliderType(entry, colliders.ColliderTypeStatic)
		transform.SetPosition(entry, geom.Vec2{
			X: object.X,
			Y: object.Y,
		})
		collision.Add(entry)
	})
}

func debugDrawEntityColliders(
	world donburi.World,
	camera engine.Camera,
	screen engine.Screen,
	collision engine.Collision,
	path *vector.Path,
	img *ebiten.Image,
	viewMatrix ebiten.GeoM) {

	worldMin := camera.ToWorld(screen.SafeArea().Min)
	worldMax := camera.ToWorld(screen.SafeArea().Max)

	entities := collision.QueryStatic(geom.AABB{
		Min: worldMin,
		Max: worldMax,
	})
	println("Static entities in view:", len(entities))

	path.Reset()

	debugPathEntityColliders(world, path, entities, viewMatrix)

	opts := &vector.DrawPathOptions{}
	opts.ColorScale.ScaleWithColor(color.RGBA{R: 200, G: 200, A: 255})
	vector.StrokePath(img, path, &vector.StrokeOptions{Width: 2}, opts)

	entities = collision.Query(geom.AABB{
		Min: worldMin,
		Max: worldMax,
	})

	path.Reset()

	debugPathEntityColliders(world, path, entities, viewMatrix)

	opts = &vector.DrawPathOptions{}
	opts.ColorScale.ScaleWithColor(color.RGBA{B: 255, A: 255})
	vector.StrokePath(img, path, &vector.StrokeOptions{Width: 2}, opts)
}

func debugPathEntityColliders(world donburi.World, path *vector.Path, entities []donburi.Entity, viewMatrix ebiten.GeoM) {
	for _, entity := range entities {
		entry := world.Entry(entity)

		position := transform.GetPosition(entry)
		collider := colliders.GetBoxCollider(entry).Translate(position)

		screenMinX, screenMinY := viewMatrix.Apply(collider.Min.X, collider.Min.Y)
		screenMaxX, screenMaxY := viewMatrix.Apply(collider.Max.X, collider.Max.Y)

		path.MoveTo(float32(screenMinX), float32(screenMinY))
		path.LineTo(float32(screenMaxX), float32(screenMinY))
		path.LineTo(float32(screenMaxX), float32(screenMaxY))
		path.LineTo(float32(screenMinX), float32(screenMaxY))
		path.Close()
	}
}
