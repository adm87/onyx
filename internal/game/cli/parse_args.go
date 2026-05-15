package cli

import "flag"

type GameArgs struct {
	Fullscreen bool
}

func ParseArgs(prog string, args []string) (*GameArgs, error) {
	gameArgs := &GameArgs{
		Fullscreen: false,
	}

	set := flag.NewFlagSet(prog, flag.ExitOnError)
	set.BoolVar(&gameArgs.Fullscreen, "fullscreen", gameArgs.Fullscreen, "run the game in fullscreen mode")

	return gameArgs, set.Parse(args)
}
