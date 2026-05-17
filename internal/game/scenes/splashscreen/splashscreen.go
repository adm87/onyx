package splashscreen

import (
	"errors"

	"github.com/adm87/onyx/internal/content"
	"github.com/adm87/onyx/pkg/components/renderer"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/yohamta/donburi"
)

const (
	CompleteExitCode engine.SceneExitCode = iota + 1
)

type sceneData struct {
	ImgEntity donburi.Entity
}

var sceneComponent = donburi.NewComponentType[sceneData]()

func getSceneData(world donburi.World) (*sceneData, bool) {
	entry, found := sceneComponent.First(world)
	return sceneComponent.Get(entry), found
}

func New(assets engine.Assets, screen engine.Screen, time engine.Time, logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		OnEnter: func(w donburi.World) error {
			return enterScene(w, assets, screen, logger)
		},
		OnExit: func(w donburi.World) error {
			return exitScene(w, assets, screen, logger)
		},
		OnUpdate: func(w donburi.World, f float64) (engine.SceneExitCode, error) {
			return CompleteExitCode, nil
		},
	}
}

func enterScene(world donburi.World, assets engine.Assets, screen engine.Screen, logger engine.Logger) error {
	if err := assets.Load(content.StaticFS(), content.Splash1920x1080Black); err != nil {
		logger.Error("failed to load splash screen asset")
		return err
	}

	cache, exists := images.GetCache(assets)
	if !exists {
		return errors.New("failed to get image cache")
	}

	img, found := cache.Get(content.Splash1920x1080Black)
	if !found {
		return errors.New("failed to get splash screen image from cache")
	}

	screen.ResizeBuffer(img.Bounds().Dx(), img.Bounds().Dy())

	image := world.Entry(world.Create(renderer.ImageRenderer...))
	renderer.Image.Get(image).Ref = img

	scene := world.Entry(world.Create(sceneComponent))
	sceneComponent.Get(scene).ImgEntity = image.Entity()
	return nil
}

func exitScene(world donburi.World, assets engine.Assets, screen engine.Screen, logger engine.Logger) error {
	assets.Unload(content.Splash1920x1080Black)
	screen.RestoreBuffer()

	sceneEntry, sceneFound := sceneComponent.First(world)
	if !sceneFound {
		logger.Warn("could not find scene component. splashscreen image may not have been removed from world. possible memory leak")
		return nil
	}

	world.Remove(sceneComponent.Get(sceneEntry).ImgEntity)
	world.Remove(sceneEntry.Entity())
	return nil
}
