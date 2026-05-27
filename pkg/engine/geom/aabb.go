package geom

type AABB struct {
	Min Vec2
	Max Vec2
}

func (a AABB) Width() float64 {
	return a.Max.X - a.Min.X
}

func (a AABB) Height() float64 {
	return a.Max.Y - a.Min.Y
}

func (a AABB) Center() Vec2 {
	return Vec2{
		X: (a.Min.X + a.Max.X) / 2,
		Y: (a.Min.Y + a.Max.Y) / 2,
	}
}

func (a AABB) Contains(point Vec2) bool {
	return point.X >= a.Min.X && point.X <= a.Max.X &&
		point.Y >= a.Min.Y && point.Y <= a.Max.Y
}

func (a AABB) Intersects(other AABB) bool {
	return a.Min.X < other.Max.X && a.Max.X > other.Min.X &&
		a.Min.Y < other.Max.Y && a.Max.Y > other.Min.Y
}

func (a AABB) Union(other AABB) AABB {
	return AABB{
		Min: Vec2{
			X: min(a.Min.X, other.Min.X),
			Y: min(a.Min.Y, other.Min.Y),
		},
		Max: Vec2{
			X: max(a.Max.X, other.Max.X),
			Y: max(a.Max.Y, other.Max.Y),
		},
	}
}
