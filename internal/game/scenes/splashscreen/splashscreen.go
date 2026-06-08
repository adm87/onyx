package splashscreen

import (
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
)

const CompleteExitCode engine.SceneExitCode = iota + 1

func New(assets engine.Assets, screen engine.Screen, world engine.World) engine.SceneState {
	var entry *donburi.Entry
	var sequence *gween.Sequence
	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			if err := assets.Load(content.EmbeddedFS(), content.EmbeddedSplash1920x1080Black); err != nil {
				return err
			}

			img, exists := images.GetImageAssets(assets, content.EmbeddedSplash1920x1080Black)
			if !exists {
				return fmt.Errorf("failed to get image asset '%s'", content.EmbeddedSplash1920x1080Black)
			}

			width, height := img.Bounds().Dx(), img.Bounds().Dy()
			halfWidth, halfHeight := float64(width)/2, float64(height)/2

			bounds := geom.AABB{
				Min: geom.Vec2{X: -halfWidth, Y: -halfHeight},
				Max: geom.Vec2{X: halfWidth, Y: halfHeight},
			}

			screen.ResizeBuffer(width, height)

			entry = images.CreateImage(ecs, content.EmbeddedSplash1920x1080Black, bounds)

			rendering.SetAnchor(entry, geom.Vec2{X: 0.5, Y: 0.5})
			rendering.SetAlpha(entry, 0)

			sequence = gween.NewSequence(
				gween.New(0, 0, 0.5, ease.Linear),
				gween.New(0, 1, 1, ease.Linear),
				gween.New(1, 1, 2, ease.Linear),
				gween.New(1, 0, 1, ease.Linear),
				gween.New(0, 0, 0.5, ease.Linear),
			)

			world.Add(entry)
			return nil
		},
		OnExit: func(ecs donburi.World) error {
			world.Remove(entry)

			assets.Unload(content.EmbeddedSplash1920x1080Black)
			screen.RestoreBuffer()

			entry = nil
			sequence = nil

			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			opacity, _, complete := sequence.Update(float32(dt))

			color := rendering.GetColor(entry)
			color.A = uint8(opacity * 255)
			rendering.SetColor(entry, color)

			if complete {
				return CompleteExitCode, nil
			}
			return engine.SceneExitNone, nil
		},
	}
}
