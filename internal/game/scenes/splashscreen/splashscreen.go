package splashscreen

import (
	"errors"

	"github.com/adm87/onyx/internal/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
)

const (
	CompleteExitCode engine.SceneExitCode = iota + 1
)

type splashscreenData struct {
	fadeSeq  *gween.Sequence
	image    *ebiten.Image
	opacity  float32
	complete bool
}

var splashscreenComponent = donburi.NewComponentType[splashscreenData]()

func New(assets engine.Assets, screen engine.Screen, time engine.Time, logger engine.Logger) *engine.SceneDefinition {
	return &engine.SceneDefinition{
		OnEnter: func(w donburi.World) error {
			return enterScene(w, assets, screen, logger)
		},
		OnExit: func(w donburi.World) error {
			return exitScene(w, assets, screen, logger)
		},
		OnUpdate: []func(world donburi.World, deltaTime float64) (engine.SceneExitCode, error){
			updateSequence,
		},
		OnDraw: []func(donburi.World, *ebiten.Image) error{
			renderSplashScreen,
		},
	}
}

func enterScene(world donburi.World, assets engine.Assets, screen engine.Screen, logger engine.Logger) error {
	if err := assets.Load(content.StaticFS(), content.Splash1920x1080Black); err != nil {
		logger.Error("Failed to load splash screen image: %v", err)
		return err
	} else if cache, exists := images.GetCache(assets); !exists {
		logger.Error("Image cache not found in assets")
		return errors.New("image cache not found in assets")
	} else if img, found := cache.Get(content.Splash1920x1080Black); !found {
		logger.Error("Failed to retrieve splash screen image from cache")
		return errors.New("splash screen image not found in cache")
	} else {
		entity := world.Create(splashscreenComponent)
		splashscreenComponent.Set(world.Entry(entity),
			&splashscreenData{
				fadeSeq: gween.NewSequence(
					gween.New(0, 0, 0.5, ease.Linear),
					gween.New(0, 1, 1.0, ease.Linear),
					gween.New(1, 1, 2.0, ease.Linear),
					gween.New(1, 0, 1.0, ease.Linear),
				),
				image: img,
			},
		)
		screen.ResizeBuffer(img.Bounds().Dx(), img.Bounds().Dy())
	}
	return nil
}

func exitScene(world donburi.World, assets engine.Assets, screen engine.Screen, logger engine.Logger) error {
	entry, found := splashscreenComponent.First(world)
	if !found {
		logger.Error("Splash screen component not found in world on exit")
		return errors.New("splash screen component not found in world on exit")
	}
	world.Remove(entry.Entity())
	assets.Unload(content.Splash1920x1080Black)
	screen.RestoreBuffer()
	return nil
}

func updateSequence(world donburi.World, deltaTime float64) (engine.SceneExitCode, error) {
	entry, found := splashscreenComponent.First(world)
	if !found {
		return engine.SceneExitCodeNone, errors.New("splash screen component not found in world")
	}

	data := splashscreenComponent.Get(entry)
	if data.complete {
		return CompleteExitCode, nil
	}

	value, _, complete := data.fadeSeq.Update(1 / 60.0)
	data.opacity = value
	data.complete = complete

	return engine.SceneExitCodeNone, nil
}

func renderSplashScreen(world donburi.World, screen *ebiten.Image) error {
	entry, found := splashscreenComponent.First(world)
	if !found {
		return errors.New("splash screen component not found in world")
	}

	data := splashscreenComponent.Get(entry)

	opt := &ebiten.DrawImageOptions{}
	opt.ColorScale.ScaleAlpha(data.opacity)

	screen.DrawImage(data.image, opt)
	return nil
}
