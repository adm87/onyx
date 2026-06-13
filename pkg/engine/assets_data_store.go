package engine

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine/file"
)

var dataFileExtensions = []file.FileExt{".json", ".xml", ".txt", ".csv", ".yaml", ".yml"}

type dataStore struct {
	store AssetStore[[]byte]
}

func newDataStore() *dataStore {
	return &dataStore{
		store: NewFileStore[[]byte](0),
	}
}

func (ds *dataStore) ImportAsset(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	ds.store.Insert(path, raw)
	return nil
}

func (ds *dataStore) DeleteAsset(path file.FilePath) bool {
	handle, exists := ds.store.GetHandle(path)
	if !exists {
		return false
	}

	_, ok := ds.store.Delete(handle)
	return ok
}

func (ds *dataStore) SupportedExtensions() []file.FileExt {
	return dataFileExtensions
}
