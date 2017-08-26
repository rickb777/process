// Wraps goroutines with the necessary synchronisation so that the caller can easily wait on completion.

package process

import "sync"

type ProcessGroup struct {
	joiner *sync.WaitGroup
}

func NewGroup() *ProcessGroup {
	return &ProcessGroup{}
}

func (pg *ProcessGroup) GoN(n int, process func()) {
	for i := 0; i < n; i++ {
		pg.Go(process)
	}
}

func (pg *ProcessGroup) Go(process func()) {
	if pg.joiner == nil {
		pg.joiner = &sync.WaitGroup{}
	}
	pg.joiner.Add(1)
	go func() {
		defer pg.joiner.Done()
		process()
	}()
}

func (pg *ProcessGroup) Join() {
	pg.joiner.Wait()
}
