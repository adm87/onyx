package transform

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type RotationData struct {
	Radians float64
	Degrees float64
}

type LocalTransformData struct {
	Matrix  ebiten.GeoM
	IsDirty bool
}

var (
	Position    = donburi.NewComponentType[geom.Vec2]()
	Scale       = donburi.NewComponentType[geom.Vec2](geom.Vec2{X: 1, Y: 1})
	Rotation    = donburi.NewComponentType[RotationData]()
	LocalMatrix = donburi.NewComponentType[LocalTransformData](LocalTransformData{IsDirty: true})
)

var Transform = []donburi.IComponentType{
	Position,
	Scale,
	Rotation,
	LocalMatrix,
}

func GetPosition(entry *donburi.Entry) (x, y float64) {
	position := Position.Get(entry)
	return position.X, position.Y
}

func SetPosition(entry *donburi.Entry, x, y float64) {
	position := Position.Get(entry)
	position.X = x
	position.Y = y
	MarkLocalDirty(entry)
}

func GetScale(entry *donburi.Entry) (x, y float64) {
	scale := Scale.Get(entry)
	return scale.X, scale.Y
}

func SetScale(entry *donburi.Entry, x, y float64) {
	scale := Scale.Get(entry)
	scale.X = x
	scale.Y = y
	MarkLocalDirty(entry)
}

func GetRediansRotation(entry *donburi.Entry) float64 {
	rotation := Rotation.Get(entry)
	return rotation.Radians
}

func SetRadiansRotation(entry *donburi.Entry, radians float64) {
	rotation := Rotation.Get(entry)
	rotation.Radians = radians
	rotation.Degrees = radians * engine.RadiansToDegrees
	MarkLocalDirty(entry)
}

func GetDegreesRotation(entry *donburi.Entry) float64 {
	rotation := Rotation.Get(entry)
	return rotation.Degrees
}

func SetDegreesRotation(entry *donburi.Entry, degrees float64) {
	rotation := Rotation.Get(entry)
	rotation.Degrees = degrees
	rotation.Radians = degrees * engine.DegreesToRadians
	MarkLocalDirty(entry)
}

func GetLocalMatrix(entry *donburi.Entry) ebiten.GeoM {
	local := LocalMatrix.Get(entry)
	if local.IsDirty {
		x, y := GetPosition(entry)
		sx, sy := GetScale(entry)
		rot := GetRediansRotation(entry)

		local.Matrix.Reset()
		local.Matrix.Scale(sx, sy)
		local.Matrix.Rotate(rot)
		local.Matrix.Translate(x, y)

		local.IsDirty = false
	}
	return local.Matrix
}

func MarkLocalDirty(entry *donburi.Entry) {
	local := LocalMatrix.Get(entry)
	local.IsDirty = true
}
