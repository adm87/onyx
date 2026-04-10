package game

import (
	onyx "github.com/adm87/onyx/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func OnStart(ctx onyx.Context) error {
	ctx.Logger().Info("game started")
	return nil
}

func OnUpdate(ctx onyx.Context) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ctx.Logger().Info("escape key pressed, exiting game")
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ctx.Logger().Info("F11 key pressed, toggling fullscreen")
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return nil
}

func OnFixedUpdate(ctx onyx.Context) error {
	return nil
}

func OnLateUpdate(ctx onyx.Context) error {
	return nil
}

func OnDraw(ctx onyx.Context) error {
	ebitenutil.DebugPrint(ctx.Screen().Image(), "Hello, Onyx!")
	return nil
}
