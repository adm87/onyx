package tiled

import (
	"encoding/xml"
	"io/fs"
	"path/filepath"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/images"
	"github.com/adm87/onyx/pkg/tiled/data"
)

var AdapterID engine.AssetAdapterID = "tiled"

type assetAdapter struct {
	logger engine.Logger

	tmxCache map[engine.FilePath]*data.Tmx
	tsxCache map[engine.FilePath]*data.Tsx

	images *images.ImageAdapter
}

func GetTmx(assets engine.Assets, path engine.FilePath) (*data.Tmx, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	tmx, exists := adapter.(*assetAdapter).tmxCache[path]
	return tmx, exists
}

func GetTsx(assets engine.Assets, path engine.FilePath) (*data.Tsx, bool) {
	adapter, found := assets.GetAdapter(AdapterID)
	if !found {
		return nil, false
	}

	tsx, exists := adapter.(*assetAdapter).tsxCache[path]
	return tsx, exists
}

func NewAdapter(images *images.ImageAdapter, logger engine.Logger) *assetAdapter {
	return &assetAdapter{
		logger:   logger,
		images:   images,
		tmxCache: make(map[engine.FilePath]*data.Tmx),
		tsxCache: make(map[engine.FilePath]*data.Tsx),
	}
}

func (a *assetAdapter) ImportAsset(fileSystem fs.FS, path engine.FilePath, raw []byte) error {
	switch path.Ext() {
	case ".tmx":
		return a.importTmx(fileSystem, path, raw)
	case ".tsx":
		return a.importTsx(fileSystem, path, raw)
	default:
		a.logger.Warn("Unsupported tiled asset file extension '%s' for asset '%s'", path.Ext(), path)
		return nil
	}
}

func (a *assetAdapter) DeleteAsset(path engine.FilePath) bool {
	return false
}

func (a *assetAdapter) SupportedExtensions() []engine.FileExt {
	return []engine.FileExt{".tmx", ".tsx"}
}

func (a *assetAdapter) importTmx(fileSystem fs.FS, path engine.FilePath, raw []byte) error {
	if _, exists := a.tmxCache[path]; exists {
		return nil
	}

	var tmx data.Tmx

	if err := xml.Unmarshal(raw, &tmx); err != nil {
		return err
	}

	tmxDir := filepath.Dir(path.String())

	tsxPaths := make([]engine.FilePath, 0, len(tmx.Tilesets))
	for i, tileset := range tmx.Tilesets {
		if tileset.Source == "" {
			continue
		}

		srcPath := a.resolvedPath(tmxDir, tileset.Source)
		tmx.Tilesets[i].Source = srcPath.String()

		tsxPaths = append(tsxPaths, srcPath)
	}

	for _, tsxPath := range tsxPaths {
		if err := a.loadTsx(fileSystem, tsxPath); err != nil {
			a.logger.Error("Failed to load tsx file '%s' referenced by tmx file '%s': %v", tsxPath, path, err)
			continue
		}
	}

	a.tmxCache[path] = &tmx
	return nil
}

func (a *assetAdapter) loadTsx(fileSystem fs.FS, path engine.FilePath) error {
	if _, exists := a.tsxCache[path]; exists {
		return nil
	}

	a.logger.Debug("Loading dependency: tsx file '%s'", path)
	raw, err := fs.ReadFile(fileSystem, path.String())

	if err != nil {
		return err
	}

	return a.importTsx(fileSystem, path, raw)
}

func (a *assetAdapter) importTsx(fileSystem fs.FS, path engine.FilePath, raw []byte) error {
	if _, exists := a.tsxCache[path]; exists {
		return nil
	}

	var tsx data.Tsx

	if err := xml.Unmarshal(raw, &tsx); err != nil {
		return err
	}

	tsxDir := filepath.Dir(path.String())

	srcPath := a.resolvedPath(tsxDir, tsx.Image.Source)
	tsx.Image.Source = srcPath.String()

	if err := a.loadTilesetImage(fileSystem, srcPath); err != nil {
		a.logger.Error("Failed to load tileset image '%s' referenced by tsx file '%s': %v", srcPath, path, err)
	}

	a.tsxCache[path] = &tsx
	return nil
}

func (a *assetAdapter) loadTilesetImage(fileSystem fs.FS, path engine.FilePath) error {
	if a.images.HasImage(path) {
		return nil
	}

	a.logger.Debug("Loading dependency: tileset image file '%s'", path)
	raw, err := fs.ReadFile(fileSystem, path.String())
	if err != nil {
		return err
	}

	return a.images.ImportAsset(fileSystem, path, raw)
}

func (a *assetAdapter) resolvedPath(directory, relativePath string) engine.FilePath {
	resolved := filepath.Join(directory, relativePath)
	return engine.FilePath(filepath.Clean(resolved))
}
