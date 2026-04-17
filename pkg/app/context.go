package app

import "context"

type Context struct {
	ctx context.Context

	cfg *Config
	lgr *Logger
	srn *Screen
}

func NewContext(ctx context.Context, cfg *Config, lgr *Logger, srn *Screen) *Context {
	return &Context{
		ctx: ctx,
		cfg: cfg,
		lgr: lgr,
		srn: srn,
	}
}

func (c *Context) WithContext(ctx context.Context) *Context {
	return &Context{
		ctx: ctx,
		cfg: c.cfg,
		lgr: c.lgr,
	}
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Config() *Config {
	return c.cfg
}

func (c *Context) Logger() *Logger {
	return c.lgr
}

func (c *Context) Screen() *Screen {
	return c.srn
}
