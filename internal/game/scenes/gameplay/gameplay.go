package gameplay

import (
	"context"
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
)

func New(assets engine.Assets, camera engine.Camera, time engine.Time) engine.SceneState {
	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			if err := assets.Load(content.AssetsFS(), content.AssetsLevelsSampleMap); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}
			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world donburi.World, screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 10, 10)
			return nil
		},
	}
}
