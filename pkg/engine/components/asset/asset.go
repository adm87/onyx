package asset

import (
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/yohamta/donburi"
)

const UnknownRef file.FilePath = "unknown"

var AssetReference = donburi.NewComponentType[file.FilePath](UnknownRef)

func NewAssetReference(world donburi.World, path file.FilePath) *donburi.Entry {
	return AddAssetReference(world.Entry(
		world.Create(
			AssetReference,
		),
	), path)
}

func AddAssetReference(entry *donburi.Entry, path file.FilePath) *donburi.Entry {
	SetAssetReference(entry, path)
	return entry
}

func GetAssetReference(entry *donburi.Entry) file.FilePath {
	if !entry.HasComponent(AssetReference) {
		return UnknownRef
	}
	return *AssetReference.Get(entry)
}

func SetAssetReference(entry *donburi.Entry, path file.FilePath) {
	donburi.Add(entry, AssetReference, &path)
}
