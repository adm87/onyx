package cli

import "flag"

type GameArgs struct {
	Fullscreen bool
}

func NewGameArgs() *GameArgs {
	return &GameArgs{
		Fullscreen: false,
	}
}

func (ga *GameArgs) Parse(prog string, args []string) error {
	set := flag.NewFlagSet(prog, flag.ContinueOnError)
	set.BoolVar(&ga.Fullscreen, "fullscreen", ga.Fullscreen, "run the game in fullscreen mode")

	if err := set.Parse(args); err != nil {
		return err
	}
	return nil
}
