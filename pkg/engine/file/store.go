package file

import "github.com/adm87/onyx/pkg/engine/storage/slotmap"

type FileStore[T any] interface {
	Insert(path FilePath, value T) uint64
	GetHandle(path FilePath) (uint64, bool)
	Get(handle uint64) (T, bool)
	Set(handle uint64, value T) (T, bool)
	Delete(handle uint64) (T, bool)
}

type fileStore[T any] struct {
	store   *slotmap.SlotMap[T]
	handles map[FilePath]uint64
}

func NewFileStore[T any](cap int) *fileStore[T] {
	return &fileStore[T]{
		store:   slotmap.New[T](cap),
		handles: make(map[FilePath]uint64),
	}
}

func (fs *fileStore[T]) Insert(path FilePath, value T) uint64 {
	handle := fs.store.Insert(value)
	fs.handles[path] = handle
	return handle
}

func (fs *fileStore[T]) GetHandle(path FilePath) (uint64, bool) {
	handle, exists := fs.handles[path]
	return handle, exists
}

func (fs *fileStore[T]) Get(handle uint64) (T, bool) {
	return fs.store.Get(handle)
}

func (fs *fileStore[T]) Set(handle uint64, value T) (T, bool) {
	return fs.store.Set(handle, value)
}

func (fs *fileStore[T]) Delete(handle uint64) (T, bool) {
	return fs.store.Delete(handle)
}
