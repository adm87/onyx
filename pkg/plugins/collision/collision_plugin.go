package collision

import "github.com/yohamta/donburi"

type CollisionPlugin struct {
	world *CollisionWorld
}

func NewCollisionPlugin() *CollisionPlugin {
	return &CollisionPlugin{
		world: NewCollisionWorld(),
	}
}

func (c *CollisionPlugin) World() *CollisionWorld {
	return c.world
}

func (c *CollisionPlugin) Add(entries ...*donburi.Entry) {
	c.world.Add(entries...)
}

func (c *CollisionPlugin) AddEntry(entry *donburi.Entry) {
	c.world.AddEntry(entry)
}

func (c *CollisionPlugin) Remove(entries ...*donburi.Entry) {
	c.world.Remove(entries...)
}

func (c *CollisionPlugin) RemoveEntry(entry *donburi.Entry) {
	c.world.RemoveEntry(entry)
}

func (c *CollisionPlugin) Update(entries ...*donburi.Entry) {
	c.world.Update(entries...)
}

func (c *CollisionPlugin) UpdateEntry(entry *donburi.Entry) {
	c.world.UpdateEntry(entry)
}
