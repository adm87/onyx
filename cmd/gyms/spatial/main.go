package main

import (
	"os"

	"github.com/adm87/onyx/internal/content"
	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
)

type test struct {
	img *ebiten.Image
}

func (t *test) Update() error {
	return nil
}

func (t *test) Draw(screen *ebiten.Image) {

}

func (t *test) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	logger := engine.NewLogger(os.Stdout)

	assets := engine.NewAssets(logger)
	assets.RegisterAdapters(images.NewImageAdapter(logger))

	if err := assets.Load(content.StaticFS(), content.Img10x10White); err != nil {
		panic(err)
	}

	cache, _ := images.GetCache(assets)
	img, _ := cache.Get(content.Img10x10White)

	if err := ebiten.RunGame(&test{img: img}); err != nil {
		panic(err)
	}

	assets.Unload(content.Img10x10White)
}
