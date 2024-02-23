package workerpool

import "sync"

type (
	slotPool struct {
		lock     sync.RWMutex
		isClosed bool
		slotChan chan slot
	}

	slot struct{}
)

func newSlotPool(size uint32) *slotPool {
	return &slotPool{
		lock:     sync.RWMutex{},
		slotChan: make(chan slot, size),
	}
}

func (sp *slotPool) reserve() bool {
	sp.lock.RLock()
	defer sp.lock.RUnlock()

	if sp.isClosed {
		return false
	}

	sp.slotChan <- slot{}

	return true
}

func (sp *slotPool) release() {
	<-sp.slotChan
}

func (sp *slotPool) close() {
	sp.lock.Lock()
	defer sp.lock.Unlock()

	sp.isClosed = true

	close(sp.slotChan)
}
