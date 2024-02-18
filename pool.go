package workerpool

import (
	"context"
	"sync"
)

type (
	worker struct{}

	Pool struct {
		wg      *sync.WaitGroup
		workers chan *worker

		keepAlive bool
	}

	Task func()
)

// NewPool returns a new pool with the size and addtional options.
func NewPool(size int, opts ...Opt) (*Pool, error) {
	poolOpts, err := parseOpts(opts...)
	if err != nil {
		return nil, err
	}

	p := &Pool{
		wg:        &sync.WaitGroup{},
		workers:   make(chan *worker, size),
		keepAlive: poolOpts.keepAlive,
	}

	if p.keepAlive {
		p.wg.Add(1)
	}

	return p, nil
}

func (p *Pool) Run(task Task) {
	w := &worker{}

	p.workers <- w // reserve a worker for the task

	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		task()

		<-p.workers // release a worker
	}()
}

func (p *Pool) Wait(ctx context.Context) error {
	done := make(chan bool, 1)
	go func() {
		p.wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
