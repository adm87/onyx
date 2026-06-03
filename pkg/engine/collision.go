package engine

type Collision interface {
}

type collision struct {
}

func newCollision(logger Logger) *collision {
	return &collision{}
}
