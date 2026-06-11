package images

import "github.com/yohamta/donburi"

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

func GetImageHandle(entry *donburi.Entry) (uint64, bool) {
	if !entry.HasComponent(ImageHandle) {
		return 0, false
	}
	return *ImageHandle.Get(entry), true
}

func SetImageHandle(entry *donburi.Entry, handle uint64) {
	donburi.Add(entry, ImageHandle, &handle)
}
