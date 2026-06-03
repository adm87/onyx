package engine

import "github.com/yohamta/donburi"

type Collision interface {
}

type collision struct {
}

func newCollision(logger Logger) *collision {
	return &collision{}
}

func (c *collision) checkCollisions(world donburi.World) error {
	return nil
}
