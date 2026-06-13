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
	Image = donburi.NewComponentType[uint64]()
	Frame = donburi.NewComponentType[int]()
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

func GetImage(entry *donburi.Entry) (uint64, bool) {
	if !entry.HasComponent(Image) {
		return 0, false
	}
	return *Image.Get(entry), true
}

func SetImageHandle(entry *donburi.Entry, handle uint64) {
	donburi.Add(entry, Image, &handle)
}

func GetFrame(entry *donburi.Entry) int {
	if !entry.HasComponent(Frame) {
		return 0
	}
	return *Frame.Get(entry)
}

func SetFrame(entry *donburi.Entry, frame int) {
	donburi.Add(entry, Frame, &frame)
}
