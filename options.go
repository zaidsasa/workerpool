package workerpool

type (
	options struct {
		keepAlive     bool
		taskQueueSize int32
	}

	Option func(*options) error
)

const defaultTaskQueueSize = 100

// WithKeepAliveOption keeps the worker pool alive even without having tasks to run (default = false).
func WithKeepAliveOption(keepAlive bool) Option {
	return func(opts *options) error {
		opts.keepAlive = keepAlive

		return nil
	}
}

// WithTaskQueueSizeOption sets the task queue size (default = 100).
func WithTaskQueueSizeOption(size int32) Option {
	return func(opts *options) error {
		opts.taskQueueSize = size

		return nil
	}
}

func parseOpts(poolOpts ...Option) (*options, error) {
	opts := &options{
		keepAlive:     false,
		taskQueueSize: defaultTaskQueueSize,
	}
	for _, opt := range poolOpts {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}
