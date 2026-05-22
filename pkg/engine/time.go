package engine

import gtime "time"

type Time interface {
	DeltaTime() float64
	FixedDeltaTime() float64
	FixedSteps() int
}

type time struct {
	deltaTime      gtime.Duration
	fixedDeltaTime gtime.Duration
	accumulator    gtime.Duration
	lastTick       gtime.Time
	fixedSteps     int
}

const (
	AccumulatorMax = gtime.Second
	DeltaTimeMax   = gtime.Second
)

func newTime(fps int) *time {
	fps = max(fps, 1)
	return &time{
		fixedDeltaTime: gtime.Second / gtime.Duration(fps),
	}
}

func (t *time) DeltaTime() float64 {
	return t.deltaTime.Seconds()
}

func (t *time) FixedDeltaTime() float64 {
	return t.fixedDeltaTime.Seconds()
}

func (t *time) FixedSteps() int {
	return t.fixedSteps
}

func (t *time) tick() {
	now := gtime.Now()

	if t.lastTick.IsZero() {
		t.lastTick = now
		return
	}

	t.deltaTime = now.Sub(t.lastTick)
	t.lastTick = now

	if t.deltaTime < 0 {
		t.deltaTime = 0
	}

	if t.deltaTime > DeltaTimeMax {
		t.deltaTime = DeltaTimeMax
	}

	t.accumulator += t.deltaTime
	t.fixedSteps = 0

	if t.accumulator > AccumulatorMax {
		t.accumulator = AccumulatorMax
	}

	for t.accumulator >= t.fixedDeltaTime {
		t.accumulator -= t.fixedDeltaTime
		t.fixedSteps++
	}
}
