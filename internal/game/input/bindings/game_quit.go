package bindings

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameQuitBinding struct {
	isActive bool
}

func NewGameQuitBinding() *GameQuitBinding {
	return &GameQuitBinding{}
}

func (b *GameQuitBinding) ID() engine.InputBindingID {
	return QuitBindingID
}

func (b *GameQuitBinding) Poll() error {
	if !b.isActive {
		return nil
	}
	if ebiten.IsKeyPressed(keyboardBindingTable[b.ID()]) {
		return ebiten.Termination
	}
	return nil
}

func (b *GameQuitBinding) SetActive(active bool) {
	b.isActive = active
}
