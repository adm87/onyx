package splashscreen

import (
	"context"
	"fmt"

	"github.com/adm87/onyx/internal/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
)

const CompleteExitCode engine.SceneExitCode = iota + 1

func New(assets engine.Assets, camera engine.Camera, time engine.Time, screen engine.Screen, logger engine.Logger) engine.SceneState {
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

			entity = world.Create(
				transform.Matrix,
				rendering.Renderer,
				rendering.Image,
			)
			entry := world.Entry(entity)

			color := rendering.GetColor(entry)
			color.A = 0

			rendering.SetImage(entry, img)
			rendering.SetAnchor(entry, geom.Vec2{X: 0.5, Y: 0.5})
			rendering.SetColor(entry, color)

			screen.ResizeBuffer(
				img.Bounds().Dx(),
				img.Bounds().Dy(),
			)
			return nil
		},
		OnExit: func(ctx context.Context, world donburi.World) error {
			screen.RestoreBuffer()
			assets.Unload(content.EmbeddedSplash1920x1080Black)
			world.Remove(entity)
			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			value, _, complete := sequence.Update(float32(time.DeltaTime()))
			entry := world.Entry(entity)

			color := rendering.GetColor(entry)
			color.A = uint8(value * 255)

			rendering.SetColor(entry, color)

			if complete {
				return CompleteExitCode, nil
			}
			return engine.SceneExitNone, nil
		},
	}
}
