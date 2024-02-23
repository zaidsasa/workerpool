package workerpool

import "sync"

func newSlotPool(size uint32) *SlotPool {
	return &SlotPool{
		lock:     &sync.RWMutex{},
		slotChan: make(chan slot, size),
	}
}

func (sp *SlotPool) reserve() bool {
	sp.lock.RLock()
	defer sp.lock.RUnlock()

	if sp.isClosed {
		return false
	}

	sp.slotChan <- slot{}

	return true
}

func (sp *SlotPool) release() {
	<-sp.slotChan
}

func (sp *SlotPool) close() {
	sp.lock.Lock()
	defer sp.lock.Unlock()

	sp.isClosed = true

	close(sp.slotChan)
}
