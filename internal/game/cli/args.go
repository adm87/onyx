package cli

import (
	"flag"
	"os"
)

type CmdArgs struct {
	Fullscreen bool
	RootDir    string
}

func ParseArgs() (CmdArgs, error) {
	args := CmdArgs{
		Fullscreen: false,
		RootDir:    ".",
	}

	set := flag.NewFlagSet("onyx-game", flag.ExitOnError)
	set.BoolVar(&args.Fullscreen, "fullscreen", args.Fullscreen, "Run the game in fullscreen mode")
	set.StringVar(&args.RootDir, "root", args.RootDir, "Root directory of the game")

	return args, set.Parse(os.Args[1:])
}
