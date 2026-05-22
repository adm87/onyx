package splashscreen

import (
	"context"
	"image/color"
	"math/rand/v2"

	"github.com/adm87/onyx/internal/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const CompleteExitCode engine.SceneExitCode = iota

func New(camera engine.Camera, time engine.Time, logger engine.Logger) engine.SceneState {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(transform.Rotation),
			filter.Not(filter.Contains(engine.CameraTag)),
		),
	)
	return engine.SceneState{
		OnEnter: func(ctx context.Context, world donburi.World) error {
			logger.Info("Entering Splash Screen Scene")

			data, err := content.EmbeddedFS().Open(content.EmbeddedImg10x10White)
			if err != nil {
				return err
			}
			defer data.Close()

			img, _, err := ebitenutil.NewImageFromReader(data)
			if err != nil {
				return err
			}

			count := 10000

			entities := world.CreateMany(count,
				transform.Matrix,
				rendering.Renderer,
				rendering.Image,
			)

			for i := range count {
				entry := world.Entry(entities[i])

				rendering.SetLayer(entry, i)
				rendering.SetImage(entry, img)
				rendering.SetAnchor(entry, geom.Vec2{X: 0.5, Y: 0.5})
				rendering.SetColor(entry, color.RGBA{
					R: uint8(rand.Float64() * 255),
					G: uint8(rand.Float64() * 255),
					B: uint8(rand.Float64() * 255),
					A: 255,
				})

				x := rand.Float64() * 1280
				y := rand.Float64() * 720
				scale := 0.5 + rand.Float64()

				transform.SetPosition(entry, geom.Vec2{X: x, Y: y})
				transform.SetScale(entry, geom.Vec2{X: scale, Y: scale})
				transform.SetRotation(entry, rand.Float64()*360)
			}

			return nil
		},
		OnExit: func(ctx context.Context, world donburi.World) error {
			logger.Info("Exiting Splash Screen Scene")
			return nil
		},
		OnUpdate: func(ctx context.Context, world donburi.World) (engine.SceneExitCode, error) {
			position := camera.Position()

			if ebiten.IsKeyPressed(ebiten.KeyW) {
				position.Y -= 200 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				position.Y += 200 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				position.X -= 200 * time.DeltaTime()
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				position.X += 200 * time.DeltaTime()
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF) {
				ebiten.SetFullscreen(!ebiten.IsFullscreen())
			}

			if position.X < 0 {
				position.X = 0
			}
			if position.Y < 0 {
				position.Y = 0
			}
			if position.X > 1280 {
				position.X = 1280
			}
			if position.Y > 720 {
				position.Y = 720
			}

			camera.SetPosition(position)

			query.Each(world, func(e *donburi.Entry) {
				rot := transform.GetRotation(e)
				rot += 100 * time.DeltaTime()
				if rot > 360 {
					rot -= 360
				}
				transform.SetRotation(e, rot)
			})
			return engine.SceneExitNone, nil
		},
	}
}
