// Package process wraps goroutines with the necessary synchronisation so that the caller
// can easily wait on completion.
//
// Following the CSP practice, these goroutines are called 'processes' - the 'P' in
// Communicating Sequential Processes. But this term is not to be confused with other
// usages, especially of processes within operating systems; the latter are unrelated.
package process

import "sync"

// ProcessGroup allows a related group of processes (i.e. goroutines) - as few as
// one - to be launched. The parent goroutine can then wait for completion of the
// entire group via `Join()`.
//
// A single parent goroutine will own each ProcessGroup. They should not be shared
// by more than one parent.
type ProcessGroup struct {
	joiner *sync.WaitGroup
}

// NewGroup creates a new empty process group. Use Go and GoN to start processes
// (i.e. goroutines) within the group.
func NewGroup() *ProcessGroup {
	return &ProcessGroup{joiner: &sync.WaitGroup{}}
}

// GoN starts n identical processes (i.e. goroutines) within this group.
// The processes are provided as zero-argument functions.
// This method can be called several times with different functions as needed.
func (pg *ProcessGroup) GoN(n int, process func()) {
	for i := 1; i <= n; i++ {
		pg.Go(process)
	}
}

// GoN1 starts n identical processes (i.e. goroutines) within this group.
// The processes are provided as one-argument functions.
// This method can be called several times with different functions as needed.
// The process argument receives the index in the sequence, starting from one.
func (pg *ProcessGroup) GoN1(n int, process func(j int)) {
	for i := 1; i <= n; i++ {
		func(j int) {
			pg.Go(func() {
				process(j)
			})
		}(i)
	}
}

// Go starts a single process (i.e. goroutine) within this group.
// This method can be called several times with different functions as needed.
func (pg *ProcessGroup) Go(process func()) {
	pg.joiner.Add(1)
	go func() {
		defer pg.joiner.Done()
		process()
	}()
}

// Join is called by the parent goroutine when it wants to sit and wait for
// every process (goroutine) in this group to have terminated. Join will therefore
// block until this condition is reached.
//
// Because the process group does not control the internal behaviour of each child
// process (goroutine), it has no means to guarantee that they will all terminated.
// So it is possible for this method to wait forever (deadlock), as a program error.
// It is up to the client code to prevent this by ensuring that all the child
// processes (goroutines) terminate cleanly.
func (pg *ProcessGroup) Join() {
	if pg.joiner != nil {
		pg.joiner.Wait()
	}
}
