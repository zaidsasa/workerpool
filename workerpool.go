package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type (
	WorkerPool struct {
		wg        *sync.WaitGroup
		workers   chan worker
		taskQueue chan Task
		keepAlive bool
	}

	Task func()

	worker struct{}
)

const minSize = 1

var ErrInvalidSize = fmt.Errorf("size must be greater than or equal to %q", minSize)

// New returns a new WorkerPool with the size and additional options.
func New(size uint32, opts ...Option) (*WorkerPool, error) {
	if size < minSize {
		return nil, ErrInvalidSize
	}

	poolOpts, err := parseOpts(opts...)
	if err != nil {
		return nil, err
	}

	workerpool := &WorkerPool{
		wg:        &sync.WaitGroup{},
		workers:   make(chan worker, size),
		taskQueue: make(chan Task, poolOpts.taskQueueSize),
		keepAlive: poolOpts.keepAlive,
	}

	go workerpool.dispatch()

	return workerpool, nil
}

// Submit queues a task for execution by the next available worker.
func (wp *WorkerPool) Submit(task Task) {
	wp.wg.Add(1)
	wp.taskQueue <- task
}

func (wp *WorkerPool) Wait(ctx context.Context) error {
	done := make(chan struct{}, 1)

	if wp.keepAlive {
		wp.wg.Add(1)
	}

	go func() {
		wp.wg.Wait()
		wp.finalize()

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
		wp.finalize()

		return fmt.Errorf("wait exited: %w", ctx.Err())
	}
}

// TODO: dispatch should recive a stopSignal to exit.
func (wp *WorkerPool) dispatch() {
	for {
		wp.workers <- worker{} // reserve a worker for the task

		go func() {
			for task := range wp.taskQueue {
				task()
				wp.wg.Done()
			}

			<-wp.workers // release a worker
		}()
	}
}

func (wp *WorkerPool) finalize() {
	// close(wp.taskQueue)
	// close(wp.workers)
}
