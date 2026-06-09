package content

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adm87/onyx/pkg/engine/file"
)

//go:embed embedded
var embedded embed.FS

func EmbeddedFS() fs.FS {
	return embedded
}

var assets fs.FS

func AssetsFS() fs.FS {
	return assets
}

const (
	AssetsAsepriteCaptainImg file.FilePath = "aseprite/captain/captain.png"
	AssetsAsepriteCaptain    file.FilePath = "aseprite/captain/captain.json"

	AssetsTiledGym01 file.FilePath = "tiled/gym01.tmx"
	AssetsTiledGym02 file.FilePath = "tiled/gym02.tmx"
	AssetsTiledGym03 file.FilePath = "tiled/gym03.tmx"
	AssetsTiledGym04 file.FilePath = "tiled/gym04.tmx"

	EmbeddedImg10x10White        file.FilePath = "embedded/images/img_10x10_white.png"
	EmbeddedSplash1920x1080Black file.FilePath = "embedded/images/splash_1920x1080_black.png"
)

func InitContentDirectories(rootDir string) {
	assets = os.DirFS(filepath.Join(rootDir, "content", "assets"))
}
