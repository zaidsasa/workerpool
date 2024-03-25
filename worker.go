package workerpool

import "sync"

type worker struct {
	wg        *sync.WaitGroup
	taskQueue chan Task
	semaphore *semaphore
}

func (w *worker) run() {
	go func() {
		for {
			task, ok := <-w.taskQueue

			if !ok {
				w.semaphore.release()

				return
			}

			task()
			w.wg.Done()
		}
	}()
}
