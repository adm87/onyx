package colliders

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/yohamta/donburi"
)

type ColliderType uint8

const (
	ColliderTypeStatic ColliderType = iota
	ColliderTypeDynamic
	ColliderTypeKinematic
)

type BodyType uint8

const (
	BodyTypeBox BodyType = iota
)

type CollisionLayer uint16

type ColliderData struct {
	Type  ColliderType
	Layer CollisionLayer
	Body  BodyType
}

var (
	Collider = donburi.NewComponentType[ColliderData](ColliderData{
		Body: BodyTypeBox,
	})
	BoxCollider = donburi.NewComponentType[geom.AABB](geom.AABB{
		Min: geom.Vec2{X: 0, Y: 0},
		Max: geom.Vec2{X: 1, Y: 1},
	})
)

func GetColliderType(entry *donburi.Entry) ColliderType {
	colliderData := Collider.Get(entry)
	if colliderData == nil {
		return ColliderTypeStatic
	}
	return colliderData.Type
}

func SetColliderType(entry *donburi.Entry, colliderType ColliderType) {
	if !entry.HasComponent(Collider) {
		entry.AddComponent(Collider)
	}
	colliderData := Collider.Get(entry)
	colliderData.Type = colliderType
}

func GetBodyType(entry *donburi.Entry) BodyType {
	colliderData := Collider.Get(entry)
	if colliderData == nil {
		return BodyTypeBox
	}
	return colliderData.Body
}

func SetBodyType(entry *donburi.Entry, bodyType BodyType) {
	if !entry.HasComponent(Collider) {
		entry.AddComponent(Collider)
	}
	colliderData := Collider.Get(entry)
	colliderData.Body = bodyType
}

func GetBoxCollider(entry *donburi.Entry) geom.AABB {
	boxCollider := BoxCollider.Get(entry)
	if boxCollider == nil {
		return geom.AABB{}
	}
	return *boxCollider
}

func SetBoxCollider(entry *donburi.Entry, aabb geom.AABB) {
	if !entry.HasComponent(BoxCollider) {
		entry.AddComponent(BoxCollider)
	}
	donburi.SetValue(entry, BoxCollider, aabb)
}
