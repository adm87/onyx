package game

import (
	"context"

	"github.com/adm87/onyx/pkg/engine"
)

func Boot() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shell := engine.NewShell(800, 600)
	shell.SetContext(ctx)

	return shell.Start()
}
