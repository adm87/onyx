package cli

import "flag"

type Args struct {
	Fullscreen bool
}

func NewArgs() *Args {
	return &Args{
		Fullscreen: false,
	}
}

func (a *Args) Parse(prog string, args []string) error {
	fs := flag.NewFlagSet(prog, flag.ContinueOnError)

	fs.BoolVar(&a.Fullscreen, "fullscreen", a.Fullscreen, "run the game in fullscreen mode")

	return fs.Parse(args)
}
