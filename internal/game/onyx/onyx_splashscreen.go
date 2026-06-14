package onyx

import (
	"fmt"
	"image/color"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
)

func (o *Onyx) SplashScreenScene() engine.SceneState {
	var splashScreenEntry *donburi.Entry
	var sequence *gween.Sequence
	var opacity float32
	var sequenceComplete bool
	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			assets := o.game.Assets()
			world := o.game.World()
			screen := o.game.Screen()

			err := assets.Load(content.EmbeddedFS(), content.EmbeddedSplash1920x1080Black)
			assert.Nil(err, fmt.Sprintf("failed to load splash screen image: %v", err))

			handle, exists := o.images.GetAssetHandle(content.EmbeddedSplash1920x1080Black)
			assert.True(exists, "failed to get handle for loaded splash screen image")

			width, height, ok := o.images.GetImageSize(handle)
			assert.True(ok, "failed to get image size for splash screen")
			screen.ResizeBuffer(width, height)

			splashScreenEntry = o.images.CreateImageEntity(ecs,
				images.WithHandle(handle),
				images.WithAnchor(0.5, 0.5),
				images.WithColor(color.RGBA{
					R: 255, G: 255, B: 255, A: 0,
				}),
			)

			transform.SetLocalBounds(splashScreenEntry, &geom.AABB{
				Min: geom.Vec2{X: -float64(width) / 2, Y: -float64(height) / 2},
				Max: geom.Vec2{X: float64(width) / 2, Y: float64(height) / 2},
			})

			world.Add(splashScreenEntry)

			sequence = gween.NewSequence(
				gween.New(0, 0, 0.25, ease.Linear),
				gween.New(0, 1, 1.5, ease.Linear),
				gween.New(1, 1, 0.5, ease.Linear),
				gween.New(1, 0, 1.5, ease.Linear),
			)

			return nil
		},
		OnExit: func(ecs donburi.World) error {
			assets := o.game.Assets()
			world := o.game.World()
			screen := o.game.Screen()

			assets.Unload(content.EmbeddedSplash1920x1080Black)
			screen.RestoreBuffer()
			world.Remove(splashScreenEntry)

			splashScreenEntry = nil
			sequence = nil

			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			if sequenceComplete {
				return SplashScreenCompleteExitCode, nil
			}

			opacity, _, sequenceComplete = sequence.Update(float32(dt))
			images.SetAlpha(splashScreenEntry, uint8(opacity*255))

			return engine.SceneExitNone, nil
		},
	}
}
