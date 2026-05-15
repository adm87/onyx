package main

import (
	"os"

	"github.com/adm87/onyx/internal/game"
)

var version string = "0.0.0-unreleased"

func main() {
	if err := game.Boot(version); err != nil {
		println("error: " + err.Error())
		os.Exit(1)
	}
}
