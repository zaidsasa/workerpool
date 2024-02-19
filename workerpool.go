package workerpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type (
	WorkerPool struct {
		wg        *sync.WaitGroup
		workers   chan *worker
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

// New returns a new WorkerPool with the size and additional options.
func New(size int, opts ...Opt) (*WorkerPool, error) {
	if size < minPoolSize {
		return nil, ErrInvalidPoolSize
	}

	poolOpts, err := parseOpts(opts...)
	if err != nil {
		return nil, err
	}

	workerpool := &WorkerPool{
		wg:        &sync.WaitGroup{},
		workers:   make(chan *worker, size),
		keepAlive: poolOpts.keepAlive,
	}

	return workerpool, nil
}

func (wp *WorkerPool) Run(task Task) {
	w := &worker{}

	wp.workers <- w // reserve a worker for the task

	wp.wg.Add(1)

	go func() {
		defer wp.wg.Done()

		task()

		<-wp.workers // release a worker
	}()
}

func (wp *WorkerPool) Wait(ctx context.Context) error {
	done := make(chan struct{}, 1)

	if wp.keepAlive {
		wp.wg.Add(1)
	}

	go func() {
		wp.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil

	case <-ctx.Done():
		if wp.keepAlive {
			wp.wg.Done()
		}

		wp.wg.Wait()

		return errors.Join(ErrContextError, ctx.Err())
	}
}
