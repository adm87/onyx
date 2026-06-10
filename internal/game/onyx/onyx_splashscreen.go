package onyx

import (
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/yohamta/donburi"
)

func (o *Onyx) SplashScreenScene() engine.SceneState {
	var splashScreenEntry *donburi.Entry

	assets := o.game.Assets()
	screen := o.game.Screen()
	world := o.game.World()

	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			err := assets.Load(content.EmbeddedFS(), content.EmbeddedSplash1920x1080Black)
			assert.Nil(err, fmt.Sprintf("failed to load splash screen image: %v", err))

			handle, exists := o.images.GetAssetHandle(content.EmbeddedSplash1920x1080Black)
			assert.True(exists, "failed to get handle for loaded splash screen image")

			splashScreenEntry = o.images.CreateImage(ecs, images.WithImageHandle(handle))
			rendering.SetAnchor(splashScreenEntry, geom.Vec2{X: 0.5, Y: 0.5})
			world.Add(splashScreenEntry)

			width, height, ok := o.images.GetImageSize(handle)
			assert.True(ok, "failed to get image size for splash screen")
			screen.ResizeBuffer(width, height)

			return nil
		},
		OnExit: func(ecs donburi.World) error {
			assets.Unload(content.EmbeddedSplash1920x1080Black)
			screen.RestoreBuffer()
			world.Remove(splashScreenEntry)

			splashScreenEntry = nil
			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			return engine.SceneExitNone, nil
		},
	}
}
