package throttle

import "sync"

type Throttle struct {
	waiting sync.WaitGroup
	ch      chan interface{}
}

func New(max int) *Throttle {
	return &Throttle{
		ch: make(chan interface{}, max),
	}
}

func (t *Throttle) Do() error {
	for {
		select {
		case t.ch <- struct{}{}:
			t.waiting.Add(1)
			return nil
		}
	}
}

func (t *Throttle) Done() {
	select {
	case <-t.ch:
	default:
		panic("Done mismatch")
	}
	t.waiting.Done()
}

func (t *Throttle) Finish() error {
	t.waiting.Wait()
	close(t.ch)
	return nil
}
