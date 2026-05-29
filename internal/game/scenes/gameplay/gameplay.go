package gameplay

import (
	"context"
	"fmt"

	"github.com/adm87/onyx-game/content"
	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/adm87/onyx-game/pkg/tiled"
	"github.com/adm87/onyx-game/pkg/tiled/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var testMaps = []engine.FilePath{
	content.AssetsLevelsGym01,
	content.AssetsLevelsGym02,
	content.AssetsLevelsGym03,
}

func New(assets engine.Assets, camera engine.Camera, screen engine.Screen, time engine.Time) engine.SceneState {
	var entity donburi.Entity

	mapIndex := 1
	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			if err := assets.Load(content.AssetsFS(), testMaps...); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}

			entry := tiled.CreateTilemap(world,
				tiled.WithTilemapRef(testMaps[mapIndex]),
			)
			entity = entry.Entity()

			tilemap, exists := tiled.GetTilemap(assets, testMaps[mapIndex])
			if !exists {
				return fmt.Errorf("tilemap asset not found: %s", testMaps[mapIndex])
			}

			camera.SetPosition(tilemap.Bounds().Center())
			camera.SetZoom(0.4)

			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			position := camera.Position()

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

			if inpututil.IsKeyJustPressed(ebiten.KeyN) {
				mapIndex = (mapIndex + 1) % len(testMaps)
				components.SetTilemapRef(world.Entry(entity), testMaps[mapIndex])
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyP) {
				mapIndex = (mapIndex - 1 + len(testMaps)) % len(testMaps)
				components.SetTilemapRef(world.Entry(entity), testMaps[mapIndex])
			}

			tilemap, exists := tiled.GetTilemap(assets, testMaps[mapIndex])
			if !exists {
				return engine.SceneExitNone, fmt.Errorf("tilemap asset not found: %s", testMaps[mapIndex])
			}

			bounds := tilemap.Bounds()

			if position.X < bounds.Min.X {
				position.X = bounds.Min.X
			}
			if position.X > bounds.Max.X {
				position.X = bounds.Max.X
			}
			if position.Y < bounds.Min.Y {
				position.Y = bounds.Min.Y
			}
			if position.Y > bounds.Max.Y {
				position.Y = bounds.Max.Y
			}

			camera.SetPosition(position)
			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world donburi.World, img *ebiten.Image, viewMatrix ebiten.GeoM) error {
			minX, minY := screen.SafeArea().Min.XY()
			ebitenutil.DebugPrintAt(img, "FPS: "+fmt.Sprintf("%.2f", ebiten.ActualFPS()), int(minX)+10, int(minY)+10)
			return nil
		},
	}
}
