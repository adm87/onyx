package data

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine/file"
)

var dataFileExtensions = []file.FileExt{".json", ".xml", ".txt", ".csv", ".yaml", ".yml"}

type DataAssets struct {
	store file.FileStore[[]byte]
}

func NewAssetAdapter() *DataAssets {
	return &DataAssets{
		store: file.NewFileStore[[]byte](0),
	}
}

func (d *DataAssets) GetData(handle uint64) ([]byte, bool) {
	return d.store.Get(handle)
}

func (d *DataAssets) GetDataHandle(path file.FilePath) (uint64, bool) {
	return d.store.GetHandle(path)
}

func (d *DataAssets) ImportAsset(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	d.store.Insert(path, raw)
	return nil
}

func (d *DataAssets) DeleteAsset(path file.FilePath) bool {
	handle, exists := d.store.GetHandle(path)
	if !exists {
		return false
	}

	_, ok := d.store.Delete(handle)
	return ok
}

func (d *DataAssets) SupportedExtensions() []file.FileExt {
	return dataFileExtensions
}
