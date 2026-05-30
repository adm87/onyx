package engine

import "github.com/adm87/onyx/pkg/engine/geom"

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func Clamp[T number](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampVec2(position, min, max geom.Vec2) geom.Vec2 {
	return geom.Vec2{
		X: Clamp(position.X, min.X, max.X),
		Y: Clamp(position.Y, min.Y, max.Y),
	}
}
