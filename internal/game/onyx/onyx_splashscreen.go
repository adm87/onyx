package onyx

import (
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/ecs/image"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
)

func (o *Onyx) SplashScreenScene() engine.SceneState {
	var cameraEntry *donburi.Entry
	var imageEntry *donburi.Entry
	var sequence *gween.Sequence
	var complete bool

	assets := o.game.Assets()
	images := o.image.Assets()
	screen := o.game.Screen()
	ecs := o.ecs
	return engine.SceneState{
		OnEnter: func() error {
			world := ecs.World()
			factory := ecs.Factory()

			if err := assets.Load(content.EmbeddedFS(), content.EmbeddedSplash1920x1080Black); err != nil {
				return err
			}

			imgHandle, exists := images.GetHandle(content.EmbeddedSplash1920x1080Black)
			if !exists {
				return fmt.Errorf("failed to get image handle for splash screen")
			}

			img, exists := images.Get(imgHandle)
			if !exists {
				return fmt.Errorf("failed to get image for splash screen")
			}

			screen.ResizeBuffer(img.Bounds().Dx(), img.Bounds().Dy())

			cameraEntry = factory.CreateCamera(world)
			imageEntry = factory.CreateImage(world,
				image.WithHandle(imgHandle),
				image.WithAnchor(0.5, 0.5),
			)
			image.SetAlpha(imageEntry, 0)
			ecs.Add(cameraEntry, imageEntry)

			sequence = gween.NewSequence(
				gween.New(0, 0, 0.1, ease.Linear),
				gween.New(0, 1, 1.0, ease.Linear),
				gween.New(1, 1, 2.0, ease.Linear),
				gween.New(1, 0, 1.0, ease.Linear),
				gween.New(0, 0, 0.1, ease.Linear),
			)

			return nil
		},
		OnUpdate: func(delta float64) (engine.SceneExitCode, error) {
			if complete {
				return SplashScreenCompleteExitCode, nil
			}

			opacity, _, seqComplete := sequence.Update(float32(delta))
			complete = seqComplete

			image.SetAlpha(imageEntry, uint8(opacity*255))
			return engine.SceneExitNone, nil
		},
		OnExit: func() error {
			ecs.Remove(cameraEntry, imageEntry)
			screen.RestoreBuffer()
			return nil
		},
	}
}
