package content

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adm87/onyx-game/pkg/engine"
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
	AssetsLevelsGym01 engine.FilePath = "levels/gym01.tmx"
	AssetsLevelsGym02 engine.FilePath = "levels/gym02.tmx"
	AssetsLevelsGym03 engine.FilePath = "levels/gym03.tmx"

	EmbeddedImg10x10White        engine.FilePath = "embedded/images/img_10x10_white.png"
	EmbeddedSplash1920x1080Black engine.FilePath = "embedded/images/splash_1920x1080_black.png"
)

func LoadDefaultContent(assets engine.Assets, logger engine.Logger) error {
	logger.Debug("Loading default content...")
	return assets.Load(embedded,
		EmbeddedImg10x10White,
	)
}

func InitContentDirectories(rootDir string) {
	assets = os.DirFS(filepath.Join(rootDir, "content", "assets"))
}
