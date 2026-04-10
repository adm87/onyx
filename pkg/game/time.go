package game

import (
	"time"
)

const (
	maxAccumulator = time.Second / 4
)

type Time interface {
	Delta32() float32
	Delta64() float64

	FixedDelta32() float32
	FixedDelta64() float64

	FixedSteps() int
}

type timeImpl struct {
	delta       time.Duration
	fixedDelta  time.Duration
	accumulator time.Duration
	last        time.Time
	steps       int
}

func (t *timeImpl) tick() {
	if t.last.IsZero() {
		t.last = time.Now()
		return
	}

	now := time.Now()
	t.delta = now.Sub(t.last)
	t.last = now

	t.accumulator += t.delta
	if t.accumulator > maxAccumulator {
		t.accumulator = maxAccumulator
	}

	t.steps = int(t.accumulator / t.fixedDelta)
	t.accumulator -= time.Duration(t.steps) * t.fixedDelta
}

func (t *timeImpl) Delta32() float32 {
	return float32(t.delta.Seconds())
}

func (t *timeImpl) Delta64() float64 {
	return t.delta.Seconds()
}

func (t *timeImpl) FixedDelta32() float32 {
	return float32(t.fixedDelta.Seconds())
}

func (t *timeImpl) FixedDelta64() float64 {
	return t.fixedDelta.Seconds()
}

func (t *timeImpl) FixedSteps() int {
	return t.steps
}
