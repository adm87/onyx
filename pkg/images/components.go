package images

import (
	"github.com/yohamta/donburi"
)

type Options struct {
	Handle uint64
	Frame  int
}

type Option func(*Options)

var (
	ImageHandle = donburi.NewComponentType[uint64]()
	ImageFrame  = donburi.NewComponentType[int]()
)

func defaultImageOptions() *Options {
	return &Options{
		Handle: 0,
		Frame:  0,
	}
}

func WithImageHandle(handle uint64) Option {
	return func(opts *Options) {
		opts.Handle = handle
	}
}

func WithImageFrame(frame int) Option {
	return func(opts *Options) {
		opts.Frame = frame
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

func GetImageFrame(entry *donburi.Entry) (int, bool) {
	if !entry.HasComponent(ImageFrame) {
		return 0, false
	}
	return *ImageFrame.Get(entry), true
}

func SetImageFrame(entry *donburi.Entry, frame int) {
	donburi.Add(entry, ImageFrame, &frame)
}
