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
	Splash1920x1080Black file.Path = "static/images/splash_1920x1080_black.png"
)
