package geom

type AABB struct {
	X, Y          float64
	Width, Height float64
}

func (a AABB) Min() (x, y float64) {
	return a.X, a.Y
}

func (a AABB) Max() (x, y float64) {
	return a.X + a.Width, a.Y + a.Height
}

func (a AABB) Center() (x, y float64) {
	return a.X + a.Width/2, a.Y + a.Height/2
}

func (a AABB) Intersects(other AABB) bool {
	return a.X < other.X+other.Width &&
		a.X+a.Width > other.X &&
		a.Y < other.Y+other.Height &&
		a.Y+a.Height > other.Y
}

func (a AABB) ContainsPoint(x, y float64) bool {
	return x >= a.X && x <= a.X+a.Width &&
		y >= a.Y && y <= a.Y+a.Height
}

func (a AABB) Contains(other AABB) bool {
	return other.X >= a.X &&
		other.Y >= a.Y &&
		other.X+other.Width <= a.X+a.Width &&
		other.Y+other.Height <= a.Y+a.Height
}
