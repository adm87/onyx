package colliders

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type BodyType uint8

const (
	BodyTypeBox BodyType = iota
)

type CollisionLayer uint16

type ColliderData struct {
	Layer CollisionLayer
	Body  BodyType
}

var (
	StaticColliderType    = donburi.NewTag("static_collider")
	DynamicColliderType   = donburi.NewTag("dynamic_collider")
	KinematicColliderType = donburi.NewTag("kinematic_collider")
)

var (
	Collider = donburi.NewComponentType[ColliderData](ColliderData{
		Body: BodyTypeBox,
	})
	BoxCollider = donburi.NewComponentType[geom.AABB](geom.AABB{
		Min: geom.Vec2{X: 0, Y: 0},
		Max: geom.Vec2{X: 1, Y: 1},
	})
)

var (
	StaticColliderQuery = donburi.NewQuery(
		filter.Contains(StaticColliderType),
	)
	DynamicColliderQuery = donburi.NewQuery(
		filter.Contains(DynamicColliderType),
	)
	KinematicColliderQuery = donburi.NewQuery(
		filter.Contains(KinematicColliderType),
	)
)

func IsStatic(entry *donburi.Entry) bool {
	return entry.HasComponent(StaticColliderType)
}

func IsDynamic(entry *donburi.Entry) bool {
	return entry.HasComponent(DynamicColliderType)
}

func IsKinematic(entry *donburi.Entry) bool {
	return entry.HasComponent(KinematicColliderType)
}

func GetBodyType(entry *donburi.Entry) BodyType {
	if !entry.HasComponent(Collider) {
		return BodyTypeBox
	}
	return Collider.Get(entry).Body
}

func SetBodyType(entry *donburi.Entry, bodyType BodyType) {
	if !entry.HasComponent(Collider) {
		entry.AddComponent(Collider)
	}
	colliderData := Collider.Get(entry)
	colliderData.Body = bodyType
}

func GetBoxCollider(entry *donburi.Entry) geom.AABB {
	if !entry.HasComponent(BoxCollider) {
		return geom.AABB{
			Min: geom.Vec2{X: 0, Y: 0},
			Max: geom.Vec2{X: 1, Y: 1},
		}
	}
	return *BoxCollider.Get(entry)
}

func SetBoxCollider(entry *donburi.Entry, aabb geom.AABB) {
	if !entry.HasComponent(BoxCollider) {
		entry.AddComponent(BoxCollider)
	}
	donburi.SetValue(entry, BoxCollider, aabb)
}
