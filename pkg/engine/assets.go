package engine

import (
	"fmt"
	"io/fs"

	"github.com/adm87/onyx/pkg/engine/file"
)

type AdapterID string

type AssetAdapter interface {
	Import(path file.Path, data []byte) error
	Delete(path file.Path)
	SupportedTypes() []file.Ext
	ID() AdapterID
}

type Assets interface {
	Load(fileSystem fs.FS, paths ...file.Path) error
	Unload(paths ...file.Path)

	RegisterAdapters(adapters ...AssetAdapter) error
	GetAdapter(id AdapterID) (AssetAdapter, bool)
}

type assets struct {
	logger Logger

	adaptersByID  map[AdapterID]AssetAdapter
	adaptersByExt map[file.Ext]AssetAdapter
}

func NewAssets(logger Logger) Assets {
	return &assets{
		logger:        logger,
		adaptersByID:  make(map[AdapterID]AssetAdapter),
		adaptersByExt: make(map[file.Ext]AssetAdapter),
	}
}

func (a *assets) Load(fileSystem fs.FS, paths ...file.Path) error {
	if len(paths) == 0 {
		a.logger.Warn("No asset paths provided to load")
		return nil
	}
	if fileSystem == nil {
		a.logger.Warn("No file system provided to load assets from")
		return nil
	}
	for _, path := range paths {
		data, err := fs.ReadFile(fileSystem, string(path))
		if err != nil {
			a.logger.Error("Failed to read asset at path: %s, error: %v", path, err)
			continue
		}

		adapter, exists := a.adaptersByExt[path.Ext()]
		if !exists {
			a.logger.Warn("No adapter registered for asset type: %s, path: %s", path.Ext(), path)
			continue
		}

		if err := adapter.Import(path, data); err != nil {
			a.logger.Error("Failed to import asset at path: %s, error: %v", path, err)
			continue
		}
		a.logger.Debug("Successfully loaded asset at path: %s", path)
	}
	return nil
}

func (a *assets) Unload(paths ...file.Path) {
	if len(paths) == 0 {
		a.logger.Warn("No asset paths provided to unload")
		return
	}

	for _, path := range paths {
		adapter, exists := a.adaptersByExt[path.Ext()]
		if !exists {
			a.logger.Warn("No adapter registered for asset type: %s, path: %s", path.Ext(), path)
			continue
		}

		adapter.Delete(path)
		a.logger.Debug("Successfully unloaded asset at path: %s", path)
	}
}

func (a *assets) RegisterAdapters(adapters ...AssetAdapter) error {
	for _, adapter := range adapters {
		id := adapter.ID()

		if _, exists := a.adaptersByID[id]; exists {
			return fmt.Errorf("adapter with ID '%s' is already registered", adapter.ID())
		}
		a.adaptersByID[id] = adapter

		for _, ext := range adapter.SupportedTypes() {
			if existingAdapter, exists := a.adaptersByExt[ext]; exists {
				return fmt.Errorf("file extension '%s' is already handled by adapter '%s'", ext, existingAdapter.ID())
			}
			a.adaptersByExt[ext] = adapter
		}
	}
	return nil
}

func (a *assets) GetAdapter(id AdapterID) (AssetAdapter, bool) {
	adapter, exists := a.adaptersByID[id]
	return adapter, exists
}
