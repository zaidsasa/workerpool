package workerpool

type (
	slotPool struct {
		slotChan chan slot
	}

	slot struct{}
)

func newSlotPool(size uint32) *slotPool {
	return &slotPool{
		slotChan: make(chan slot, size),
	}
}

func (sp *slotPool) acquire() {
	sp.slotChan <- slot{}
}

func (sp *slotPool) release() {
	<-sp.slotChan
}

func (sp *slotPool) close() {
	close(sp.slotChan)
}
