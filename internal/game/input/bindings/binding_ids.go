package bindings

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	QuitBindingID engine.InputBindingID = iota
	FullscreenToggleBindingID
)

var keyboardBindingTable = map[engine.InputBindingID]ebiten.Key{
	QuitBindingID:             ebiten.KeyEscape,
	FullscreenToggleBindingID: ebiten.KeyF11,
}
