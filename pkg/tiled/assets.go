package tiled

import (
	"io/fs"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/storage/slotmap"
)

type assetAdapter struct {
	tmxStore *slotmap.SlotMap[*Tmx]
	tsxStore *slotmap.SlotMap[*Tsx]

	tmxHandles map[file.FilePath]uint64
	tsxHandles map[file.FilePath]uint64
}

func newAssetsAdapter(assets engine.Assets) *assetAdapter {
	return &assetAdapter{
		tmxStore:   slotmap.New[*Tmx](0),
		tsxStore:   slotmap.New[*Tsx](0),
		tmxHandles: make(map[file.FilePath]uint64),
		tsxHandles: make(map[file.FilePath]uint64),
	}
}

func (a *assetAdapter) ImportAsset(fileSystem fs.FS, path file.FilePath, raw []byte) error {
	return nil
}

func (a *assetAdapter) DeleteAsset(path file.FilePath) bool {
	return true
}

func (a *assetAdapter) SupportedExtensions() []file.FileExt {
	return []file.FileExt{"tmx", "tsx"}
}
