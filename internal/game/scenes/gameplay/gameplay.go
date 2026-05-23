package gameplay

import (
	"context"
	"fmt"
	"image/color"
	"math/rand"

	"github.com/adm87/onyx/internal/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

func New(assets engine.Assets, camera engine.Camera, time engine.Time) engine.SceneState {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(transform.Rotation),
			filter.Not(
				filter.Contains(engine.CameraTag),
			),
		),
	)
	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			img, found := images.GetImage(assets, content.EmbeddedImg10x10White)
			if !found {
				return fmt.Errorf("failed to load image: %s", content.EmbeddedImg10x10White)
			}
			for range 10000 {
				scale := 0.5 + rand.Float64()
				images.NewEntity(world,
					images.WithImage(img),
					images.WithPosition(
						1280*rand.Float64(),
						720*rand.Float64(),
					),
					images.WithAnchor(0.5, 0.5),
					images.WithScale(scale, scale),
					images.WithColor(
						color.RGBA{
							R: uint8(rand.Intn(256)),
							G: uint8(rand.Intn(256)),
							B: uint8(rand.Intn(256)),
							A: 255,
						},
					),
				)
			}
			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			position := camera.Position()

			if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
				position.X -= 5
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
				position.X += 5
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
				position.Y -= 5
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
				position.Y += 5
			}

			camera.SetPosition(position)

			query.Each(world, func(entry *donburi.Entry) {
				transform.Rotate(entry, 100*time.DeltaTime())
			})

			return engine.SceneExitNone, nil
		},
		OnRender: func(ctx context.Context, world donburi.World, screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 10, 10)
			return nil
		},
	}
}
