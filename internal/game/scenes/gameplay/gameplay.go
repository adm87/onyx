package gameplay

import (
	"context"
	"fmt"

	"github.com/adm87/onyx-game/content"
	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/adm87/onyx-game/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

func New(assets engine.Assets, camera engine.Camera, screen engine.Screen, time engine.Time) engine.SceneState {
	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			if err := assets.Load(content.AssetsFS(), content.AssetsLevelsSampleMap); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}
			tiled.CreateTilemap(world,
				tiled.WithTilemapRef(content.AssetsLevelsSampleMap),
			)
			camera.SetZoom(0.2)
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

			camera.SetPosition(position)
			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world donburi.World, screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
			ebitenutil.DebugPrint(screen, "FPS: "+fmt.Sprintf("%.2f", ebiten.ActualFPS()))
			return nil
		},
	}
}
