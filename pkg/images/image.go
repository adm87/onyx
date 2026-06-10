package images

import (
	"github.com/yohamta/donburi"
)

type ImageOptions struct {
	Handle uint64
}

type ImageOption func(*ImageOptions)

var ImageHandle = donburi.NewComponentType[uint64]()

func defaultImageOptions() *ImageOptions {
	return &ImageOptions{
		Handle: 0,
	}
}

func WithImageHandle(handle uint64) ImageOption {
	return func(opts *ImageOptions) {
		opts.Handle = handle
	}
}

func GetImageHandle(entry *donburi.Entry) uint64 {
	if !entry.HasComponent(ImageHandle) {
		return 0
	}
	return *ImageHandle.Get(entry)
}

func SetImageHandle(entry *donburi.Entry, handle uint64) {
	donburi.Add(entry, ImageHandle, &handle)
}
