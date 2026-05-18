package quadtree

type Quadtree interface {
}

type node struct {
}

func New() Quadtree {
	return &node{}
}
