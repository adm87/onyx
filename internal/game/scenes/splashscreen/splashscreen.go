package splashscreen

import (
	"context"
	"fmt"

	"github.com/adm87/onyx/internal/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/images"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
)

const CompleteExitCode engine.SceneExitCode = iota + 1

func updateSplashscreen(world donburi.World, entity donburi.Entity, value float32) {
	entry := world.Entry(entity)

	color := rendering.GetColor(entry)
	color.A = uint8(value * 255)

	rendering.SetColor(entry, color)
}

func New(assets engine.Assets, time engine.Time, screen engine.Screen) engine.SceneState {
	var entity donburi.Entity
	sequence := gween.NewSequence(
		gween.New(0, 0, 0.5, ease.Linear),
		gween.New(0, 1, 1, ease.Linear),
		gween.New(1, 1, 2, ease.Linear),
		gween.New(1, 0, 1, ease.Linear),
		gween.New(0, 0, 0.5, ease.Linear),
	)
	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			if err := assets.Load(content.EmbeddedFS(), content.EmbeddedSplash1920x1080Black); err != nil {
				return err
			}

			img, exists := images.GetImage(assets, content.EmbeddedSplash1920x1080Black)
			if !exists {
				return fmt.Errorf("failed to get image asset '%s'", content.EmbeddedSplash1920x1080Black)
			}

			screen.ResizeBuffer(img.Bounds().Dx(), img.Bounds().Dy())

			entity = images.NewEntity(world,
				images.WithImage(img),
				images.WithAnchor(0.5, 0.5),
				images.WithAlpha(0),
			)

			return nil
		},
		OnExit: func(ctx context.Context, world donburi.World) error {
			world.Remove(entity)

			assets.Unload(content.EmbeddedSplash1920x1080Black)
			screen.RestoreBuffer()

			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			value, _, complete := sequence.Update(float32(time.DeltaTime()))
			updateSplashscreen(world, entity, value)

			if complete {
				return CompleteExitCode, nil
			}

			return engine.SceneExitNone, nil
		},
	}
}
