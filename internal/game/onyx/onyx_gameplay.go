package onyx

import (
	"fmt"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/debug"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/plugins/ecs/camera"
	"github.com/adm87/onyx/pkg/plugins/ecs/tiled"
	"github.com/adm87/onyx/pkg/plugins/ecs/transform"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var gameplayManifest = []file.FilePath{
	content.AssetsAsepriteCaptainImg,
	content.AssetsAsepriteCaptainJson,
	content.AssetsTiledGym04,
}

func (o *Onyx) GameplayScene() engine.SceneState {
	var cameraEntry *donburi.Entry
	var tilemapEntry *donburi.Entry

	assets := o.game.Assets()
	screen := o.game.Screen()

	tiledAssets := o.tiled.Assets()

	ecs := o.ecs
	return engine.SceneState{
		OnEnter: func() error {
			if err := assets.Load(content.AssetsFS(), gameplayManifest...); err != nil {
				return err
			}

			tmxHandle, exists := tiledAssets.GetTmxHandle(content.AssetsTiledGym04)
			if !exists {
				return fmt.Errorf("failed to get TMX handle for gym")
			}

			tilemap := tiledAssets.BuildTilemap(tmxHandle)
			tilemapCenter := tilemap.Bounds().Center()

			cameraEntry = ecs.Factory().CreateCamera(ecs.World())
			camera.SetZoom(cameraEntry, 0.25)

			transform.SetPosition(cameraEntry, tilemapCenter.X, tilemapCenter.Y)

			tilemapEntry = ecs.Factory().CreateTilemap(ecs.World(),
				tiled.WithTilemapHandle(tmxHandle),
			)
			ecs.Add(tilemapEntry)

			return nil
		},
		OnRender: func(target *ebiten.Image) error {
			debug.DrawTransformBounds(ecs, cameraEntry, target, screen.SafeArea())
			return nil
		},
	}
}
