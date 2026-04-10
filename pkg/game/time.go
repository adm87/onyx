package game

type Time interface {
	Delta32() float32
	Delta64() float64

	FixedDelta32() float32
	FixedDelta64() float64

	FixedSteps() int
}

type timeImpl struct {
	delta32      float32
	delta64      float64
	fixedDelta32 float32
	fixedDelta64 float64
	fixedSteps   int
}

func (t *timeImpl) tick() {

}

func (t *timeImpl) Delta32() float32 {
	return t.delta32
}

func (t *timeImpl) Delta64() float64 {
	return t.delta64
}

func (t *timeImpl) FixedDelta32() float32 {
	return t.fixedDelta32
}

func (t *timeImpl) FixedDelta64() float64 {
	return t.fixedDelta64
}

func (t *timeImpl) FixedSteps() int {
	return t.fixedSteps
}
