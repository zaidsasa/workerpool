package workerpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type (
	Pool struct {
		wg      *sync.WaitGroup
		workers chan *worker

		keepAlive bool
	}

	Task func()

	worker struct{}
)

const minPoolSize = 1

var (
	ErrInvalidPoolSize = fmt.Errorf("pool size must be bigger or equal to %d", minPoolSize)
	ErrContextError    = errors.New("pool context error")
)

// NewPool returns a new pool with the size and additional options.
func NewPool(size int, opts ...Opt) (*Pool, error) {
	if size < minPoolSize {
		return nil, ErrInvalidPoolSize
	}

	poolOpts, err := parseOpts(opts...)
	if err != nil {
		return nil, err
	}

	pool := &Pool{
		wg:        &sync.WaitGroup{},
		workers:   make(chan *worker, size),
		keepAlive: poolOpts.keepAlive,
	}
	if pool.keepAlive {
		pool.wg.Add(1)
	}

	return pool, nil
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
		return errors.Join(ErrContextError, ctx.Err())
	}
}
