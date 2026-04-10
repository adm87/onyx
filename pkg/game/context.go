package game

import "context"

type Context interface {
	Context() context.Context
	Logger() Logger
	Model() Model
	Screen() Screen
	Time() Time
}

type ctxImpl struct {
	ctx context.Context
	lgr Logger
	mdl Model

	srn *screenImpl
	tm  *timeImpl
}

func (c *ctxImpl) Context() context.Context {
	return c.ctx
}

func (c *ctxImpl) Logger() Logger {
	return c.lgr
}

func (c *ctxImpl) Model() Model {
	return c.mdl
}

func (c *ctxImpl) Screen() Screen {
	return c.srn
}

func (c *ctxImpl) Time() Time {
	return c.tm
}
