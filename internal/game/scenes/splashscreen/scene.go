package splashscreen

import (
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/ecs"
	"github.com/adm87/onyx/pkg/plugins/ecs/camera"
	"github.com/adm87/onyx/pkg/plugins/ecs/transform"
	"github.com/adm87/onyx/pkg/plugins/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
)

const (
	SplashScreenCompleteExitCode engine.SceneExitCode = iota + 1
)

type Scene struct {
	game engine.Game

	imgEntry *donburi.Entry
	camEntry *donburi.Entry

	sequence    *gween.Sequence
	seqComplete bool
}

func NewScene(game engine.Game) *Scene {
	return &Scene{
		game: game,
		sequence: gween.NewSequence(
			gween.New(0, 0, 0.25, ease.Linear),
			gween.New(0, 1, 1.0, ease.InCubic),
			gween.New(1, 1, 1.5, ease.Linear),
			gween.New(1, 0, 1.0, ease.OutCubic),
			gween.New(0, 0, 0.25, ease.Linear),
		),
	}
}

func (s *Scene) Enter() error {
	assets := s.game.Assets()
	if err := assets.Load(content.EmbeddedFS(), content.EmbeddedSplash1920x1080Black); err != nil {
		return err
	}

	imagePlugin := engine.GetPlugin[images.ImagePlugin](s.game, images.PluginID())
	imageAssets := imagePlugin.Assets()

	handle, exists := imageAssets.GetHandle(content.EmbeddedSplash1920x1080Black)
	if !exists {
		return fmt.Errorf("failed to get handle for splash screen image")
	}

	width, height, _ := imageAssets.GetImageSize(handle)

	screen := s.game.Screen()
	screen.ResizeBuffer(width, height)

	ecsPlugin := engine.GetPlugin[ecs.ECSPlugin](s.game, ecs.PluginID())

	s.imgEntry = imagePlugin.CreateImage(ecsPlugin.World(),
		images.WithHandle(handle),
		images.WithAnchor(0.5, 0.5),
	)

	s.camEntry = transform.NewTransform(ecsPlugin.World())
	s.camEntry.AddComponent(camera.MainCamera)

	ecsPlugin.Add(s.imgEntry, s.camEntry)
	return nil
}

func (s *Scene) Exit() error {
	ecsPlugin := engine.GetPlugin[ecs.ECSPlugin](s.game, ecs.PluginID())
	ecsPlugin.Remove(s.imgEntry, s.camEntry)

	screen := s.game.Screen()
	screen.RestoreBuffer()

	assets := s.game.Assets()
	assets.Unload(content.EmbeddedSplash1920x1080Black)
	return nil
}

func (s *Scene) Update(dt float64) (engine.SceneExitCode, error) {
	if s.seqComplete {
		return SplashScreenCompleteExitCode, nil
	}

	opacity, _, seqComplete := s.sequence.Update(float32(dt))
	s.seqComplete = seqComplete

	images.SetAlpha(s.imgEntry, uint8(opacity*255))
	camera.SetZoom(s.camEntry, float64(0.95+0.05*opacity))

	return engine.SceneExitNone, nil
}

func (s *Scene) FixedUpdate(dt float64) error {
	return nil
}

func (s *Scene) LateUpdate(dt float64) error {
	return nil
}

func (s *Scene) Render(target *ebiten.Image) error {
	return nil
}
