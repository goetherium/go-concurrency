package semaphore

import (
	"errors"
	"time"
)

type Semaphore struct {
	tickets chan struct{}
	timeout time.Duration
}

func New(tickets int, timeout time.Duration) *Semaphore {
	return &Semaphore{
		timeout: timeout,
		tickets: make(chan struct{}, tickets),
	}
}

func (s Semaphore) Acquire() error {
	select {
	case s.tickets <- struct{}{}:
		return nil
	case <-time.After(s.timeout):
		return errors.New("timeout")
	}
}

func (s Semaphore) Release() {
	<-s.tickets
}
