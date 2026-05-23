package content

import (
	"embed"
	"io/fs"

	"github.com/adm87/onyx/pkg/engine"
)

//go:embed embedded
var embedded embed.FS

func EmbeddedFS() fs.FS {
	return embedded
}

const (
	EmbeddedImg10x10White        engine.FilePath = "embedded/images/img_10x10_white.png"
	EmbeddedSplash1920x1080Black engine.FilePath = "embedded/images/splash_1920x1080_black.png"
)
