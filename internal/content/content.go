package content

import (
	"embed"
	"io/fs"
)

//go:embed embedded
var embedded embed.FS

func EmbeddedFS() fs.FS {
	return embedded
}

const (
	EmbeddedImg10x10White = "embedded/images/img_10x10_white.png"
)
