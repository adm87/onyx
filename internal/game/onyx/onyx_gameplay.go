package onyx

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (o *Onyx) GameplayScene() engine.SceneState {
	return engine.SceneState{
		OnEnter: func() error {
			return nil
		},
		OnRender: func(target *ebiten.Image) error {
			ebitenutil.DebugPrint(target, "Gameplay Scene")
			return nil
		},
	}
}
