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

	var tilemap *tiled.Tilemap
	var entity donburi.Entity

	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			if err := assets.Load(content.AssetsFS(), tilemapRef); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}
			tiled.CreateTilemap(world,
				tiled.WithTilemapRef(tilemapRef),
			)

			tm, exists := tiled.GetTilemap(assets, tilemapRef)
			if !exists {
				return fmt.Errorf("tilemap asset not found: %s", tilemapRef)
			}
			tilemap = tm

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

			camera.SetPosition(position)
			transform.SetPosition(entry, position)

			collision.Update(entry)
			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world donburi.World, img *ebiten.Image, viewMatrix ebiten.GeoM) error {
			partitioning.DebugDrawSpatialHash(img, collision.Partitioning(), screen.SafeArea(), viewMatrix)

			worldMin := camera.ToWorld(screen.SafeArea().Min)
			worldMax := camera.ToWorld(screen.SafeArea().Max)

			entities := collision.Query(geom.AABB{
				Min: worldMin,
				Max: worldMax,
			})

			for _, entity := range entities {
				entry := world.Entry(entity)

				position := transform.GetPosition(entry)
				collider := colliders.GetBoxCollider(entry).Translate(position)

				screenMinX, screenMinY := viewMatrix.Apply(collider.Min.X, collider.Min.Y)
				screenMaxX, screenMaxY := viewMatrix.Apply(collider.Max.X, collider.Max.Y)

				vector.StrokeRect(img,
					float32(screenMinX),
					float32(screenMinY),
					float32(screenMaxX-screenMinX),
					float32(screenMaxY-screenMinY),
					2,
					color.RGBA{0, 255, 0, 100},
					false,
				)
			}
			return nil
		},
	}
}
