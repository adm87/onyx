package tiled

import (
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/tiled/components"
	"github.com/yohamta/donburi"
)

type TilemapOptions struct {
	TilemapRef engine.FilePath
}

type TilemapOption func(*TilemapOptions)

func defaultTilemapOptions() *TilemapOptions {
	return &TilemapOptions{
		TilemapRef: "",
	}
}

func CreateTilemap(world donburi.World, opts ...TilemapOption) *donburi.Entry {
	entity := world.Create(
		components.Tilemap,
	)
	entry := world.Entry(entity)

	options := defaultTilemapOptions()
	for _, opt := range opts {
		opt(options)
	}

	components.SetTilemapRef(entry, options.TilemapRef)

	return entry
}

func WithTilemapRef(ref engine.FilePath) TilemapOption {
	return func(opts *TilemapOptions) {
		opts.TilemapRef = ref
	}
}
