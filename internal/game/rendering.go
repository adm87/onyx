package game

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
)

func addRenderingSystems(onyx engine.Game) {
	renderer := onyx.Renderer()

	renderer.AddRenderingSystem(images.NewRenderingSystem())
}
