package workerpool

import "sync"

type worker struct {
	wg        *sync.WaitGroup
	taskQueue chan Task
	slots     *slotPool
}

func (w *worker) run() {
	go func() {
		for {
			task, ok := <-w.taskQueue

			if !ok {
				w.slots.release()

				return
			}

			task()
			w.wg.Done()
		}
	}()
}
