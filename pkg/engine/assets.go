package engine

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
)

type AdapterID string

type AssetAdapter interface {
	ImportAsset(assets Assets, fileSystem fs.FS, path file.FilePath, raw []byte) error
	DeleteAsset(path file.FilePath) bool
	SupportedExtensions() []file.FileExt
}

type Assets interface {
	Load(fileSystem fs.FS, paths ...file.FilePath) error
	Unload(paths ...file.FilePath)

	AddAssetAdapter(adapter AssetAdapter) uint64
	GetAdapter(id uint64) (AssetAdapter, bool)

	GetDataHandle(path file.FilePath) (uint64, bool)
	GetData(handle uint64) ([]byte, bool)
}

type assets struct {
	logger *logger
	store  *slotmap.SlotMap[AssetAdapter]

	adaptersByExt map[file.FileExt]uint64
	dataAssets    *DataAssets
}

func newAssets(logger *logger) *assets {
	a := &assets{
		logger:        logger,
		store:         slotmap.New[AssetAdapter](0),
		adaptersByExt: make(map[file.FileExt]uint64),
		dataAssets:    NewAssetAdapter(),
	}
	a.AddAssetAdapter(a.dataAssets)
	return a
}

func (a *assets) GetDataHandle(path file.FilePath) (uint64, bool) {
	return a.dataAssets.GetDataHandle(path)
}

func (a *assets) GetData(handle uint64) ([]byte, bool) {
	return a.dataAssets.GetData(handle)
}

func (a *assets) Load(fileSystem fs.FS, paths ...file.FilePath) error {
	for _, path := range paths {
		raw, err := fs.ReadFile(fileSystem, path.String())
		if err != nil {
			a.logger.Error("Failed to read asset file '%s': %v", path, err)
			continue
		}

		handle, exists := a.adaptersByExt[path.Ext()]
		if !exists {
			a.logger.Error("No asset adapter found for file extension '%s', skipping asset '%s'", path.Ext(), path)
			continue
		}

		adapter, exists := a.store.Get(handle)
		if !exists {
			a.logger.Error("Asset adapter with handle %d not found for asset '%s'", handle, path)
			continue
		}

		if err := adapter.ImportAsset(a, fileSystem, path, raw); err != nil {
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

		handle, exists := a.adaptersByExt[ext]
		if !exists {
			a.logger.Warn("No asset adapter found for file extension '%s', skipping unload of asset '%s'", ext, path)
			continue
		}

		adapter, exists := a.store.Get(handle)
		if !exists {
			a.logger.Warn("Asset adapter with handle %d not found for asset '%s', skipping unload", handle, path)
			continue
		}

		if adapter.DeleteAsset(path) {
			a.logger.Debug("Successfully unloaded asset '%s'", path)
		} else {
			a.logger.Warn("Failed to unload asset '%s': asset not found in adapter", path)
		}
	}
}

func (a *assets) AddAssetAdapter(adapter AssetAdapter) uint64 {
	handle := a.store.Insert(adapter)

	for _, ext := range adapter.SupportedExtensions() {
		if _, exists := a.adaptersByExt[ext]; exists {
			a.logger.Warn("Asset adapter for extension '%s' already exists, skipping", ext)
			continue
		}

		a.adaptersByExt[ext] = handle
	}

	return handle
}

func (a *assets) GetAdapter(id uint64) (AssetAdapter, bool) {
	adapter, exists := a.store.Get(id)
	return adapter, exists
}
