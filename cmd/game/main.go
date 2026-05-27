package main

import (
	"os"

	"github.com/adm87/onyx-game/internal/game"
)

func main() {
	if err := game.Boot(); err != nil {
		println("error: " + err.Error())
		os.Exit(1)
	}
}
