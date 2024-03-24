package workerpool

type (
	options struct {
		keepAlive bool
	}

	Option func(*options) error
)

// WithKeepAliveOption keeps the worker pool alive even without having tasks to run (default = false).
func WithKeepAliveOption(keepAlive bool) Option {
	return func(opts *options) error {
		opts.keepAlive = keepAlive

		return nil
	}
}

func parseOpts(poolOpts ...Option) (*options, error) {
	opts := &options{
		keepAlive: false,
	}
	for _, opt := range poolOpts {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}
