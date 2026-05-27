package gameplay

import (
	"context"
	"fmt"

	"github.com/adm87/onyx-game/content"
	"github.com/adm87/onyx-game/pkg/engine"
	"github.com/yohamta/donburi"
)

func New(assets engine.Assets, camera engine.Camera, screen engine.Screen, time engine.Time) engine.SceneState {
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
	}
}
