package tiled

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/images"
)

type TiledAssetAdapter struct {
	tmxCache map[file.FilePath]*Tmx
	tsxCache map[file.FilePath]*Tsx
	tilemaps map[file.FilePath]*Tilemap

	images *images.ImageAssetAdapter
}

func NewTiledAssetAdapter(images *images.ImageAssetAdapter) *TiledAssetAdapter {
	return &TiledAssetAdapter{
		images:   images,
		tmxCache: make(map[file.FilePath]*Tmx),
		tsxCache: make(map[file.FilePath]*Tsx),
		tilemaps: make(map[file.FilePath]*Tilemap),
	}
}

func (a *TiledAssetAdapter) ImportAsset(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	switch path.Ext() {
	case ".tmx":
		return a.importTmx(fileSystem, path, raw)
	case ".tsx":
		return a.importTsx(fileSystem, path, raw)
	default:
		return fmt.Errorf("unsupported file extension '%s' for tiled asset '%s'", path.Ext(), path)
	}
}

func (a *TiledAssetAdapter) DeleteAsset(path file.FilePath) bool {
	delete(a.tmxCache, path)
	delete(a.tsxCache, path)
	delete(a.tilemaps, path)
	return true
}

func (a *TiledAssetAdapter) SupportedExtensions() []file.FileExt {
	return []file.FileExt{".tmx", ".tsx"}
}

func (a *TiledAssetAdapter) importTmx(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	if _, exists := a.tmxCache[path]; exists {
		return nil
	}

	var tmx Tmx

	if err := xml.Unmarshal(raw, &tmx); err != nil {
		return err
	}

	tmxDir := filepath.Dir(path.String())

	tsxPaths := make([]file.FilePath, 0, len(tmx.Tilesets))
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
			return fmt.Errorf("failed to load tsx file '%s' referenced by tmx file '%s': %w", tsxPath, path, err)
		}
	}

	a.tmxCache[path] = &tmx

	tilemap, err := buildTilemap(&tmx)
	if err != nil {
		return fmt.Errorf("failed to decode tilemap for tmx file '%s': %w", path, err)
	}

	a.tilemaps[path] = tilemap
	return nil
}

func (a *TiledAssetAdapter) loadTsx(fileSystem fs.FS, path file.FilePath) error {
	if _, exists := a.tsxCache[path]; exists {
		return nil
	}

	raw, err := fs.ReadFile(fileSystem, path.String())

	if err != nil {
		return err
	}

	return a.importTsx(fileSystem, path, raw)
}

func (a *TiledAssetAdapter) importTsx(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	if _, exists := a.tsxCache[path]; exists {
		return nil
	}

	var tsx Tsx

	if err := xml.Unmarshal(raw, &tsx); err != nil {
		return err
	}

	tsxDir := filepath.Dir(path.String())

	srcPath := a.resolvedPath(tsxDir, tsx.Image.Source)
	tsx.Image.Source = srcPath.String()

	if err := a.loadTilesetImage(fileSystem, srcPath); err != nil {
		return fmt.Errorf("failed to load tileset image '%s' referenced by tsx file '%s': %w", srcPath, path, err)
	}

	a.tsxCache[path] = &tsx
	return nil
}

func (a *TiledAssetAdapter) loadTilesetImage(fileSystem fs.FS, path file.FilePath) error {
	if a.images.HasImage(path) {
		return nil
	}

	raw, err := fs.ReadFile(fileSystem, path.String())
	if err != nil {
		return err
	}

	return a.images.ImportAsset(fileSystem, path, raw)
}

func (a *TiledAssetAdapter) resolvedPath(directory, relativePath string) file.FilePath {
	resolved := filepath.Join(directory, relativePath)
	return file.FilePath(filepath.Clean(resolved))
}
