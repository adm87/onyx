package ecs

import "github.com/yohamta/donburi"

type ECS interface {
	NewEntity(components ...donburi.IComponentType) *donburi.Entry
}

type donburiESC struct {
	world donburi.World
}

func NewDonburiESC() ECS {
	return &donburiESC{
		world: donburi.NewWorld(),
	}
}

func (d *donburiESC) NewEntity(components ...donburi.IComponentType) *donburi.Entry {
	return d.world.Entry(d.world.Create(components...))
}
