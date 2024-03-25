package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type (
	WorkerPool struct {
		wg        *sync.WaitGroup
		stopChan  chan struct{}
		semaphore *semaphore
		taskQueue chan Task
		keepAlive bool
	}

	Task func()
)

const (
	minSize         = 1
	taskQueueFactor = 10
)

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
		semaphore: newSemaphore(size),
		stopChan:  make(chan struct{}, 1),
		taskQueue: make(chan Task, size*taskQueueFactor),
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
	defer wp.finalize()

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
	case <-ctx.Done():
		if wp.keepAlive {
			wp.wg.Done()
		}

		wp.wg.Wait()

		return fmt.Errorf("wait exited: %w", ctx.Err())
	}

	return nil
}

func (wp *WorkerPool) dispatch() {
	for {
		select {
		case <-wp.stopChan:
			wp.semaphore.close()

			return
		default:
			wp.semaphore.acquire()

			worker := worker{
				wg:        wp.wg,
				taskQueue: wp.taskQueue,
				semaphore: wp.semaphore,
			}

			worker.run()
		}
	}
}

func (wp *WorkerPool) finalize() {
	close(wp.taskQueue)

	wp.stopChan <- struct{}{}
}
