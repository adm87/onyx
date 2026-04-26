package game

import (
	"context"
	"image/color"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type onyx struct {
	ctx    context.Context
	logger *engine.Logger
	assets *engine.Assets
	input  *engine.Input
	screen *engine.Screen
}

func newGame(
	ctx context.Context,
	cfg *engine.Config,
	logger *engine.Logger,
	input *engine.Input,
	assets *engine.Assets,
	screen *engine.Screen) *onyx {
	return &onyx{
		ctx:    ctx,
		logger: logger,
		assets: assets,
		input:  input,
		screen: screen,
	}
}

func (o *onyx) Update() error {
	select {
	case <-o.ctx.Done():
		return o.ctx.Err()
	default:
		if err := o.input.Poll(); err != nil {
			return err
		}
		return nil
	}
}

func (o *onyx) Draw(screen *ebiten.Image) {
	select {
	case <-o.ctx.Done():
		return
	default:
		buffer := o.screen.Buffer()
		buffer.Clear()

		safeMinX, safeMinY := o.screen.SafeArea().Min()
		safeMaxX, safeMaxY := o.screen.SafeArea().Max()

		safeMinX += 10
		safeMinY += 10
		safeMaxX -= 10
		safeMaxY -= 10

		// top left
		vector.FillRect(buffer, float32(safeMinX), float32(safeMinY), 100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255}, true)

		// top right
		vector.FillRect(buffer, float32(safeMaxX-100), float32(safeMinY), 100, 100, color.RGBA{R: 0, G: 255, B: 0, A: 255}, true)

		// bottom left
		vector.FillRect(buffer, float32(safeMinX), float32(safeMaxY-100), 100, 100, color.RGBA{R: 0, G: 0, B: 255, A: 255}, true)

		// bottom right
		vector.FillRect(buffer, float32(safeMaxX-100), float32(safeMaxY-100), 100, 100, color.RGBA{R: 255, G: 255, B: 0, A: 255}, true)

		screen.DrawImage(buffer, o.screen.Options())
	}
}

func (o *onyx) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return o.screen.Layout(outsideWidth, outsideHeight)
}
