package bindings

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type FullscreenToggleBinding struct {
	isActive bool
}

func NewFullscreenToggleBinding() *FullscreenToggleBinding {
	return &FullscreenToggleBinding{}
}

func (b *FullscreenToggleBinding) ID() engine.InputBindingID {
	return FullscreenToggleBindingID
}

func (b *FullscreenToggleBinding) Poll() error {
	if !b.isActive {
		return nil
	}
	if inpututil.IsKeyJustPressed(keyboardBindingTable[b.ID()]) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return nil
}

func (b *FullscreenToggleBinding) SetActive(active bool) {
	b.isActive = active
}
