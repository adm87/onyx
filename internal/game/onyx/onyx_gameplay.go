package onyx

import (
	"image/color"

	"github.com/adm87/onyx/content"
	"github.com/adm87/onyx/pkg/ecs/camera"
	"github.com/adm87/onyx/pkg/ecs/transform"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/plugins/tiled"
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
	return engine.SceneState{
		OnEnter: func() error {
			assets := o.game.Assets()
			screen := o.game.Screen()
			tiledAssets := o.tiled.Assets()

			screen.SetBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255})

			if err := assets.Load(content.AssetsFS(), gameplayManifest...); err != nil {
				return err
			}

			tmxHandle, found := tiledAssets.GetTmxHandle(content.AssetsTiledGym04)
			if !found {
				return engine.ErrAssetNotFound{Path: content.AssetsTiledGym04.String()}
			}

			tilemap := tiledAssets.BuildTilemap(tmxHandle)
			tilemapCenter := tilemap.Bounds().Center()

			tilemapEntry = o.tiled.CreateTilemap(o.ecs.World(),
				tiled.WithTilemapHandle(tmxHandle),
			)

			cameraEntry = transform.NewTransform(o.ecs.World())
			cameraEntry.AddComponent(camera.MainCamera)

			transform.SetPosition(cameraEntry, tilemapCenter.X, tilemapCenter.Y)
			camera.SetZoom(cameraEntry, 0.25)

			o.ecs.Add(cameraEntry, tilemapEntry)
			return nil
		},
		OnExit: func() error {
			o.ecs.Remove(cameraEntry, tilemapEntry)
			return nil
		},
		OnRender: func(target *ebiten.Image) error {

			return nil
		},
	}
}
