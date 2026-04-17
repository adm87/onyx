package game

import "github.com/adm87/onyx/pkg/app"

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
		return nil
	}
}
