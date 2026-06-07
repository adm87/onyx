package engine

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine/file"
)

type AdapterID string

type AssetAdapter interface {
	ImportAsset(fileSystem fs.FS, path file.FilePath, raw []byte) error
	DeleteAsset(path file.FilePath) bool
	SupportedExtensions() []file.FileExt
}

type Assets interface {
	Load(fileSystem fs.FS, paths ...file.FilePath) error
	Unload(paths ...file.FilePath)

	AddAssetAdapter(id AdapterID, adapter AssetAdapter)
	GetAdapter(id AdapterID) (AssetAdapter, bool)

	GetData(path file.FilePath) ([]byte, bool)
}

type assets struct {
	logger *logger

	adaptersByID  map[AdapterID]AssetAdapter
	adaptersByExt map[file.FileExt]AssetAdapter

	data map[file.FilePath][]byte
}

func newAssets(logger *logger) *assets {
	return &assets{
		logger:        logger,
		adaptersByID:  make(map[AdapterID]AssetAdapter),
		adaptersByExt: make(map[file.FileExt]AssetAdapter),
		data:          make(map[file.FilePath][]byte),
	}
}

func (a *assets) Load(fileSystem fs.FS, paths ...file.FilePath) error {
	for _, path := range paths {
		raw, err := fs.ReadFile(fileSystem, path.String())
		if err != nil {
			a.logger.Error("Failed to read asset file '%s': %v", path, err)
			continue
		}

		adapter, exists := a.adaptersByExt[path.Ext()]
		if !exists {
			a.data[path] = raw
			continue
		}

		if err := adapter.ImportAsset(fileSystem, path, raw); err != nil {
			a.logger.Error("Failed to import asset '%s': %v", path, err)
			continue
		}

		a.logger.Debug("Successfully loaded asset '%s'", path)
	}
	return nil
}

func (a *assets) Unload(paths ...file.FilePath) {
	for _, path := range paths {
		ext := path.Ext()

		adapter, exists := a.adaptersByExt[ext]
		if !exists {
			a.logger.Warn("No asset adapter found for file extension '%s', skipping unload of asset '%s'", ext, path)
			continue
		}

		if adapter.DeleteAsset(path) {
			a.logger.Debug("Successfully unloaded asset '%s'", path)
		} else {
			a.logger.Warn("Failed to unload asset '%s': asset not found in adapter", path)
		}
	}
}

func (a *assets) AddAssetAdapter(id AdapterID, adapter AssetAdapter) {
	if _, exists := a.adaptersByID[id]; exists {
		a.logger.Warn("Asset adapter with ID '%s' already exists, skipping", id)
		return
	}

	a.adaptersByID[id] = adapter

	for _, ext := range adapter.SupportedExtensions() {
		if _, exists := a.adaptersByExt[ext]; exists {
			a.logger.Warn("Asset adapter for extension '%s' already exists, skipping", ext)
			continue
		}

		a.adaptersByExt[ext] = adapter
	}
}

func (a *assets) GetAdapter(id AdapterID) (AssetAdapter, bool) {
	adapter, exists := a.adaptersByID[id]
	return adapter, exists
}

func (a *assets) GetData(path file.FilePath) ([]byte, bool) {
	data, exists := a.data[path]
	return data, exists
}
