package workerpool

type (
	opts struct {
		keepAlive bool
	}

	Opt func(*opts) error
)

// WithKeepAliveOption keeps the worker pool alive even without having tasks to run.
func WithKeepAliveOption() Opt {
	return func(opts *opts) error {
		opts.keepAlive = true

		return nil
	}
}

func parseOpts(poolOpts ...Opt) (*opts, error) {
	opts := &opts{
		keepAlive: false,
	}
	for _, opt := range poolOpts {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}
