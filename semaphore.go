package workerpool

type (
	semaphore struct {
		stack chan struct{}
	}
)

func newSemaphore(size uint32) *semaphore {
	return &semaphore{
		stack: make(chan struct{}, size),
	}
}

func (s *semaphore) acquire() {
	s.stack <- struct{}{}
}

func (s *semaphore) release() {
	<-s.stack
}

func (s *semaphore) close() {
	close(s.stack)
}
