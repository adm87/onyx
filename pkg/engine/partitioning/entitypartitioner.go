package partitioning

import "github.com/yohamta/donburi"

type EntityPartitioner interface {
	Accepts(entry *donburi.Entry) bool
	Add(entry *donburi.Entry)
	Remove(entry *donburi.Entry)
	Update(entry *donburi.Entry)
	Query(ecs donburi.World, region any, callback func(*donburi.Entry))
	QueryInto(ecs donburi.World, region any, result []*donburi.Entry) []*donburi.Entry
}
