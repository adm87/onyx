package onyx

import (
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
)

var gameplayManifest = []file.FilePath{
	content.AssetsTiledGym04,
}

func (o *Onyx) GameplayScene() engine.SceneState {
	var tilemapEntry *donburi.Entry
	var tilemap *tiled.Tilemap
	var tilemapHandle uint64
	var err error

	assets := o.game.Assets()
	camera := o.game.Camera()
	screen := o.game.Screen()
	world := o.game.World()

	return engine.SceneState{
		OnEnter: func(ecs donburi.World) error {
			err = assets.Load(content.AssetsFS(), gameplayManifest...)
			assert.Nil(err, fmt.Sprintf("failed to load gameplay assets: %v", err))

			tmxHandle, exists := o.tiled.GetTmxHandle(content.AssetsTiledGym04)
			assert.True(exists, "failed to get handle for tiled map")

			tilemap, tilemapHandle, err = o.tiled.BuildTilemap(tmxHandle)
			assert.Nil(err, fmt.Sprintf("failed to parse tiled map: %v", err))

			tilemapEntry = o.tiled.CreateTilemap(ecs, tiled.WithTilemapHandle(tilemapHandle))
			world.Add(tilemapEntry)

			camera.SetPosition(tilemap.Bounds().Center())
			camera.SetZoom(0.25)

			return nil
		},
		OnUpdate: func(ecs donburi.World, dt float64) (engine.SceneExitCode, error) {
			return engine.SceneExitNone, nil
		},
		OnFixedUpdate: func(ecs donburi.World, dt float64) error {
			return nil
		},
		OnLateUpdate: func(ecs donburi.World, dt float64) error {
			return nil
		},
		OnRender: func(ecs donburi.World, img *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error {
			min := screen.SafeArea().Min
			ebitenutil.DebugPrintAt(img, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), int(min.X), int(min.Y))
			return nil
		},
	}
}
