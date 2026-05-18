package transform

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RotationData struct {
	Radians float64
	Degrees float64
}

type LocalMatrixData struct {
	Matrix  ebiten.GeoM
	IsDirty bool
}

var (
	Position    = donburi.NewComponentType[geom.Vec2]()
	Scale       = donburi.NewComponentType[geom.Vec2](geom.Vec2{X: 1, Y: 1})
	Rotation    = donburi.NewComponentType[RotationData]()
	LocalMatrix = donburi.NewComponentType[LocalMatrixData](LocalMatrixData{IsDirty: true})
)

func GetPosition(entry *donburi.Entry) (x, y float64) {
	if !entry.HasComponent(Position) {
		panic("Entity does not have a Position component")
	}
	vec2 := Position.Get(entry)
	return vec2.X, vec2.Y
}

func SetPosition(entry *donburi.Entry, x, y float64) {
	if !entry.HasComponent(Position) {
		entry.AddComponent(Position)
	}
	vec2 := Position.Get(entry)
	vec2.X = x
	vec2.Y = y
	MarkDirty(entry)
}

func GetScale(entry *donburi.Entry) (x, y float64) {
	if !entry.HasComponent(Scale) {
		panic("Entity does not have a Scale component")
	}
	vec2 := Scale.Get(entry)
	return vec2.X, vec2.Y
}

func SetScale(entry *donburi.Entry, x, y float64) {
	if !entry.HasComponent(Scale) {
		entry.AddComponent(Scale)
	}
	vec2 := Scale.Get(entry)
	vec2.X = x
	vec2.Y = y
	MarkDirty(entry)
}

func GetRotationRadians(entry *donburi.Entry) float64 {
	if !entry.HasComponent(Rotation) {
		panic("Entity does not have a Rotation component")
	}
	return Rotation.Get(entry).Radians
}

func SetRotationRadians(entry *donburi.Entry, radians float64) {
	if !entry.HasComponent(Rotation) {
		entry.AddComponent(Rotation)
	}
	rotation := Rotation.Get(entry)
	rotation.Radians = radians
	rotation.Degrees = radians * 57.29577951
	MarkDirty(entry)
}

func GetRotationDegrees(entry *donburi.Entry) float64 {
	if !entry.HasComponent(Rotation) {
		panic("Entity does not have a Rotation component")
	}
	return Rotation.Get(entry).Degrees
}

func SetRotationDegrees(entry *donburi.Entry, degrees float64) {
	if !entry.HasComponent(Rotation) {
		entry.AddComponent(Rotation)
	}
	rotation := Rotation.Get(entry)
	rotation.Degrees = degrees
	rotation.Radians = degrees * 0.0174532925
	MarkDirty(entry)
}

func GetLocalMatrix(entry *donburi.Entry) ebiten.GeoM {
	if !entry.HasComponent(LocalMatrix) {
		entry.AddComponent(LocalMatrix)
	}
	localMatrix := LocalMatrix.Get(entry)
	if localMatrix.IsDirty {
		posX, posY := GetPosition(entry)
		scaleX, scaleY := GetScale(entry)
		rotation := GetRotationRadians(entry)

		localMatrix.Matrix.Reset()
		localMatrix.Matrix.Scale(scaleX, scaleY)
		localMatrix.Matrix.Rotate(rotation)
		localMatrix.Matrix.Translate(posX, posY)

		localMatrix.IsDirty = false
	}
	return localMatrix.Matrix
}

func MarkDirty(entry *donburi.Entry) {
	if !entry.HasComponent(LocalMatrix) {
		entry.AddComponent(LocalMatrix)
	}
	LocalMatrix.Get(entry).IsDirty = true
}
