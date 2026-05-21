package engine

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Shell interface {
	Start() error

	Context() context.Context
	SetContext(ctx context.Context)
}

type shell struct {
	ctx context.Context
}

func NewShell(width, height int) Shell {
	return &shell{
		ctx: context.Background(),
	}
}

func (s *shell) Context() context.Context {
	return s.ctx
}

func (s *shell) SetContext(ctx context.Context) {
	if s.ctx != nil {
		return
	}
	s.ctx = ctx
}

func (s *shell) Start() error {
	return ebiten.RunGame(s)
}

func (s *shell) Update() error {
	return nil
}

func (s *shell) Draw(screen *ebiten.Image) {
	select {
	case <-s.ctx.Done():
		return
	default:
		ebitenutil.DebugPrint(screen, "Hello, World!")
	}
}

func (s *shell) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
