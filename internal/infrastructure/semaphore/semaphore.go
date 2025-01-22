package semaphore

type Semaphore struct {
	tickets int
	ch      chan struct{}
}

func New(tickets int) *Semaphore {
	return &Semaphore{
		tickets: tickets,
		ch:      make(chan struct{}, tickets),
	}
}

func (s Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s Semaphore) Release() {
	<-s.ch
}
