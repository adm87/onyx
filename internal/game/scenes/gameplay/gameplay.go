package gameplay

import (
	"context"
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled"
	"github.com/yohamta/donburi"
)

func New(assets engine.Assets, camera engine.Camera, screen engine.Screen, time engine.Time) engine.SceneState {
	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			if err := assets.Load(content.AssetsFS(), content.AssetsLevelsSampleMap); err != nil {
				return fmt.Errorf("failed to load level asset: %w", err)
			}

			tmx, exists := tiled.GetTmx(assets, content.AssetsLevelsSampleMap)
			if !exists {
				return fmt.Errorf("failed to get tmx asset '%s'", content.AssetsLevelsSampleMap)
			}

			tsx, exists := tiled.GetTsx(assets, engine.FilePath(tmx.Tilesets[0].Source))
			if !exists {
				return fmt.Errorf("failed to get tsx asset '%s'", tmx.Tilesets[0].Source)
			}

			img, exists := images.GetImage(assets, engine.FilePath(tsx.Image.Source))
			if !exists {
				return fmt.Errorf("failed to get image asset '%s'", tsx.Image.Source)
			}

			images.CreateImage(world,
				images.WithRef(img),
				images.WithScale(4, 4),
				images.WithAnchor(0.5, 0.5),
			)

			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			return engine.SceneExitNone, nil
		},
	}
}
