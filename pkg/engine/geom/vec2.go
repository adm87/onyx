package geom

import "math"

type Vec2 struct {
	X, Y float64
}

func (v Vec2) XY() (float64, float64) {
	return v.X, v.Y
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vec2) Mul(scalar float64) Vec2 {
	return Vec2{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v Vec2) Div(scalar float64) Vec2 {
	return Vec2{
		X: v.X / scalar,
		Y: v.Y / scalar,
	}
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vec2) Normalize() Vec2 {
	length := v.Length()
	if length == 0 {
		return Vec2{0, 0}
	}
	return v.Div(length)
}

func (v Vec2) Dot(other Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec2) Cross(other Vec2) float64 {
	return v.X*other.Y - v.Y*other.X
}

func (v Vec2) Rotate(angle float64) Vec2 {
	rad := angle * math.Pi / 180
	cos := math.Cos(rad)
	sin := math.Sin(rad)
	return Vec2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}
