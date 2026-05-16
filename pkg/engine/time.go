package engine

import gotime "time"

type Time interface {
	DeltaTime() float64
	FixedDeltaTime() float64
	FixedSteps() int
	Tick()
}

type time struct {
	deltaTime  gotime.Duration
	fixedTime  gotime.Duration
	fixedSteps int
}

func NewTime() Time {
	return &time{
		deltaTime:  0,
		fixedTime:  0,
		fixedSteps: 0,
	}
}

func (t *time) DeltaTime() float64 {
	return t.deltaTime.Seconds()
}

func (t *time) FixedDeltaTime() float64 {
	return t.fixedTime.Seconds()
}

func (t *time) FixedSteps() int {
	return t.fixedSteps
}

func (t *time) Tick() {
	now := gotime.Now()
	t.deltaTime = now.Sub(gotime.Unix(0, 0))
	t.fixedTime = gotime.Duration(16_666_667) // 60 FPS
}
