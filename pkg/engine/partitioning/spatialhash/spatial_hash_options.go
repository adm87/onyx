package spatialhash

type Options struct {
	Capacity    int
	Padding     [4]int
	Resolutions []float64
}

type Option func(*Options)

func defaultOptions() Options {
	return Options{
		Capacity:    0,
		Padding:     [4]int{0, 0, 0, 0},
		Resolutions: []float64{64, 128, 256},
	}
}

func WithCapacity(capacity int) Option {
	return func(opts *Options) {
		opts.Capacity = capacity
	}
}

func WithPadding(left, top, right, bottom int) Option {
	return func(opts *Options) {
		opts.Padding = [4]int{left, top, right, bottom}
	}
}

func WithResolutions(resolutions ...float64) Option {
	return func(opts *Options) {
		opts.Resolutions = resolutions
	}
}
