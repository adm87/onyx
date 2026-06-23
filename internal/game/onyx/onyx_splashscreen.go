package onyx

import (
	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/ecs/camera"
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
)

func (o *Onyx) SplashScreenScene() engine.SceneState {
	var imageEntry *donburi.Entry
	var cameraEntry *donburi.Entry
	var sequence *gween.Sequence
	var complete bool
	return engine.SceneState{
		OnEnter: func() error {
			assets := o.game.Assets()
			screen := o.game.Screen()
			imageAssets := o.images.Assets()

			if err := assets.Load(content.EmbeddedFS(), content.EmbeddedSplash1920x1080Black); err != nil {
				return err
			}

			imgHandle, found := imageAssets.GetHandle(content.EmbeddedSplash1920x1080Black)
			if !found {
				return engine.ErrAssetNotFound{Path: content.EmbeddedSplash1920x1080Black.String()}
			}

			width, height, _ := imageAssets.GetImageSize(imgHandle)
			screen.ResizeBuffer(width, height)

			imageEntry = o.images.CreateImage(o.ecs.World(),
				images.WithHandle(imgHandle),
				images.WithAnchor(0.5, 0.5),
			)
			images.SetAlpha(imageEntry, 0)

			cameraEntry = transform.NewTransform(o.ecs.World())
			cameraEntry.AddComponent(camera.MainCamera)

			sequence = gween.NewSequence(
				gween.New(0, 0, 0.25, ease.Linear),
				gween.New(0, 1, 1.0, ease.InCubic),
				gween.New(1, 1, 1.5, ease.Linear),
				gween.New(1, 0, 1.0, ease.OutCubic),
				gween.New(0, 0, 0.25, ease.Linear),
			)

			o.Add(imageEntry, cameraEntry)
			return nil
		},
		OnUpdate: func(delta float64) (engine.SceneExitCode, error) {
			if complete {
				return SplashScreenCompleteExitCode, nil
			}

			opacity, _, seqComplete := sequence.Update(float32(delta))
			complete = seqComplete

			images.SetAlpha(imageEntry, uint8(opacity*255))
			camera.SetZoom(cameraEntry, float64(0.95+0.05*opacity))

			return engine.SceneExitNone, nil
		},
		OnExit: func() error {
			o.game.Screen().RestoreBuffer()
			o.Remove(imageEntry, cameraEntry)
			return nil
		},
	}
}
