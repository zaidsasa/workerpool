package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type (
	WorkerPool struct {
		wg        *sync.WaitGroup
		slots     *SlotPool
		taskQueue chan Task
		keepAlive bool
	}

	Task func()

	slot struct{}

	SlotPool struct {
		lock     *sync.RWMutex
		isClosed bool
		slotChan chan slot
	}
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
		slots:     newSlotPool(size),
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
		if ok := wp.slots.reserve(); !ok {
			return
		}

		go worker(wp.slots, wp.taskQueue, wp.wg)
	}
}

func (wp *WorkerPool) finalize() {
	close(wp.taskQueue)
	wp.slots.close()
}

func worker(slots *SlotPool, taskQueue chan Task, waitGroup *sync.WaitGroup) {
	for {
		task, ok := <-taskQueue

		if !ok {
			slots.release()

			return
		}

		task()
		waitGroup.Done()
	}
}
