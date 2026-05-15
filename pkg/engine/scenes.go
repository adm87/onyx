package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Scenes interface {
	Update() error
	FixedUpdate() error
	LateUpdate() error
	Render(*ebiten.Image) error
}

type scenes struct {
	world  donburi.World
	logger Logger
}

func NewScenes(world donburi.World, logger Logger) Scenes {
	return &scenes{
		world:  world,
		logger: logger,
	}
}

func (s *scenes) Update() error {
	return nil
}

func (s *scenes) FixedUpdate() error {
	return nil
}

func (s *scenes) LateUpdate() error {
	return nil
}

func (s *scenes) Render(screen *ebiten.Image) error {
	return nil
}
