package game

import (
	"image/color"

	"github.com/adm87/onyx/pkg/app"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
}

func New() *Game {
	return &Game{}
}

func (g *Game) Startup(ctx *app.Context) error {
	return nil
}

func (g *Game) Shutdown(ctx *app.Context) error {
	return nil
}

func (g *Game) Update(ctx *app.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (g *Game) Draw(ctx *app.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		minx, miny := ctx.Screen().SafeMin()
		maxx, maxy := ctx.Screen().SafeMax()

		// top left
		vector.FillRect(ctx.Screen().Buffer(), float32(minx), float32(miny), 100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255}, true)

		// top right
		vector.FillRect(ctx.Screen().Buffer(), float32(maxx-100), float32(miny), 100, 100, color.RGBA{R: 0, G: 255, B: 0, A: 255}, true)

		// bottom left
		vector.FillRect(ctx.Screen().Buffer(), float32(minx), float32(maxy-100), 100, 100, color.RGBA{R: 0, G: 0, B: 255, A: 255}, true)

		// bottom right
		vector.FillRect(ctx.Screen().Buffer(), float32(maxx-100), float32(maxy-100), 100, 100, color.RGBA{R: 255, G: 255, B: 0, A: 255}, true)

		return nil
	}
}
