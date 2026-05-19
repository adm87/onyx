package content

import (
	"embed"
	"io/fs"

	"github.com/adm87/onyx/pkg/engine/file"
)

//go:embed static
var staticFS embed.FS

func StaticFS() fs.FS {
	return staticFS
}

const (
	Img10x10White        file.Path = "static/images/img_10x10_white.png"
	Splash1920x1080Black file.Path = "static/images/splash_1920x1080_black.png"
)
