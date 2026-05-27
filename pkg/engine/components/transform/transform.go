package transform

import (
	"github.com/adm87/onyx-game/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type MatrixData struct {
	GeoM    ebiten.GeoM
	IsDirty bool
}

var (
	Position = donburi.NewComponentType[geom.Vec2]()
	Scale    = donburi.NewComponentType[geom.Vec2](geom.Vec2{X: 1, Y: 1})
	Rotation = donburi.NewComponentType[float64]()
	Matrix   = donburi.NewComponentType[MatrixData](MatrixData{IsDirty: true})
)

func GetPosition(entry *donburi.Entry) geom.Vec2 {
	if !entry.HasComponent(Position) {
		return geom.Vec2{X: 0, Y: 0}
	}
	return *Position.Get(entry)
}

func SetPosition(entry *donburi.Entry, pos geom.Vec2) {
	if !entry.HasComponent(Position) {
		entry.AddComponent(Position)
	}
	donburi.SetValue(entry, Position, pos)
	MarkDirty(entry)
}

func GetScale(entry *donburi.Entry) geom.Vec2 {
	if !entry.HasComponent(Scale) {
		return geom.Vec2{X: 1, Y: 1}
	}
	return *Scale.Get(entry)
}

func SetScale(entry *donburi.Entry, scale geom.Vec2) {
	if !entry.HasComponent(Scale) {
		entry.AddComponent(Scale)
	}
	donburi.SetValue(entry, Scale, scale)
	MarkDirty(entry)
}

func GetRotation(entry *donburi.Entry) float64 {
	if !entry.HasComponent(Rotation) {
		return 0
	}
	return *Rotation.Get(entry)
}

func SetRotation(entry *donburi.Entry, rot float64) {
	if !entry.HasComponent(Rotation) {
		entry.AddComponent(Rotation)
	}
	donburi.SetValue(entry, Rotation, rot)
	MarkDirty(entry)
}

func Rotate(entry *donburi.Entry, delta float64) float64 {
	current := GetRotation(entry)
	newRot := current + delta
	if newRot >= 360 {
		newRot -= 360
	} else if newRot < 0 {
		newRot += 360
	}
	SetRotation(entry, newRot)
	return newRot
}

func GetMatrix(entry *donburi.Entry) ebiten.GeoM {
	if !entry.HasComponent(Matrix) {
		return ebiten.GeoM{}
	}

	matrix := Matrix.Get(entry)
	if matrix.IsDirty {
		pos := GetPosition(entry)
		scale := GetScale(entry)
		rot := GetRotation(entry)

		matrix.GeoM.Reset()
		matrix.GeoM.Scale(scale.X, scale.Y)
		matrix.GeoM.Rotate(rot * 0.0174533) // Convert degrees to radians
		matrix.GeoM.Translate(pos.X, pos.Y)

		matrix.IsDirty = false
	}
	return matrix.GeoM
}

func IsDirty(entry *donburi.Entry) bool {
	if !entry.HasComponent(Matrix) {
		return false
	}
	return Matrix.Get(entry).IsDirty
}

func MarkDirty(entry *donburi.Entry) {
	if !entry.HasComponent(Matrix) {
		entry.AddComponent(Matrix)
	}
	donburi.SetValue(entry, Matrix, MatrixData{IsDirty: true})
}
