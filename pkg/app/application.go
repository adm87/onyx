package app

type Application interface {
	Startup(ctx *Context) error
	Shutdown(ctx *Context) error
	Update(ctx *Context) error
	Draw(ctx *Context) error
}
