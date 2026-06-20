package onyx

import (
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/ecs/image"
)

func (o *Onyx) SplashScreenScene() engine.SceneState {
	assets := o.game.Assets()
	images := o.image.Assets()
	screen := o.game.Screen()
	world := o.ecs.World()
	factory := o.ecs.Factory()
	return engine.SceneState{
		OnEnter: func() error {
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

			factory.CreateCamera(world)
			factory.CreateImage(world,
				image.WithHandle(imgHandle),
				image.WithAnchor(0.5, 0.5),
			)
			return nil
		},
		OnExit: func() error {
			screen.RestoreBuffer()
			return nil
		},
	}
}
