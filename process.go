// Package process wraps goroutines with the necessary synchronisation so that the caller
// can easily wait on completion.
//
// Following the CSP practice, these goroutines are called 'processes' - the 'P' in
// Communicating Sequential Processes. But this term is not to be confused with other
// usages, especially of processes within operating systems; the latter are unrelated.
package process

import (
	"sync"

	"github.com/joeshaw/multierror"
)

// ProcessGroup allows a related group of processes (i.e. goroutines) - zero or
// more - to be launched. The parent goroutine can then wait for completion of the
// entire group via `Join()`.
//
// A single parent goroutine will own each ProcessGroup. They should not be shared
// by more than one parent.
type ProcessGroup struct {
	joiner *sync.WaitGroup
	errs   multierror.Errors
	µ      sync.Mutex
}

// NewGroup creates a new empty process group. Use Go and GoN to start processes
// (i.e. goroutines) within the group.
func NewGroup() *ProcessGroup {
	return &ProcessGroup{joiner: &sync.WaitGroup{}}
}

//-------------------------------------------------------------------------------------------------
// Methods that use a process returning nothing.

// Go starts a single process (i.e. goroutine) within this group
// using a zero-argument function.
// This method can be called multiple times with different functions as needed.
// Use Join or JoinE to wait for all the processes to terminate.
func (pg *ProcessGroup) Go(process func()) {
	pg.GoN(1, process)
}

// GoN starts n identical processes (i.e. goroutines) within this group
// using a zero-argument function.
// This method can be called several times with different functions as needed.
// Use Join or JoinE to wait for all the processes to terminate.
func (pg *ProcessGroup) GoN(n int, process func()) {
	pg.joiner.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer pg.joiner.Done()
			process()
		}()
	}
}

// GoN0 starts n identical processes (i.e. goroutines) within this group
// using a one-argument function.
// This method can be called multiple times with different functions as needed.
// The process argument receives the index in the sequence, starting from zero.
// Use Join or JoinE to wait for all the processes to terminate.
func (pg *ProcessGroup) GoN0(n int, process func(j int)) {
	pg.goN(0, n, process)
}

// GoN1 starts n identical processes (i.e. goroutines) within this group
// using a one-argument function.
// This method can be called multiple times with different functions as needed.
// The process argument receives the index in the sequence, starting from one.
func (pg *ProcessGroup) GoN1(n int, process func(j int)) {
	pg.goN(1, n+1, process)
}

func (pg *ProcessGroup) goN(from, n int, process func(j int)) {
	pg.joiner.Add(n - from)
	for i := from; i < n; i++ {
		go func(j int) {
			defer pg.joiner.Done()
			process(j)
		}(i)
	}
}

//-------------------------------------------------------------------------------------------------
// Methods that use a process returning error.

// GoE starts a single process (i.e. goroutine) within this group
// using a zero-argument function returning an optional error.
// This method can be called multiple times with different functions as needed.
// Use Join or JoinE to wait for all the processes to terminate.
func (pg *ProcessGroup) GoE(process func() error) {
	pg.GoNE(1, process)
}

// GoNE starts n identical processes (i.e. goroutines) within this group
// using a zero-argument function returning an optional error.
// This method can be called several times with different functions as needed.
// Use Join or JoinE to wait for all the processes to terminate.
func (pg *ProcessGroup) GoNE(n int, process func() error) {
	pg.goNE(0, n,
		func(j int) error {
			return process()
		})
}

// GoN0E starts n identical processes (i.e. goroutines) within this group
// using a one-argument function returning an optional error.
// This method can be called multiple times with different functions as needed.
// The process argument receives the index in the sequence, starting from zero.
// Use Join or JoinE to wait for all the processes to terminate.
func (pg *ProcessGroup) GoN0E(n int, process func(j int) error) {
	pg.goNE(0, n, process)
}

// GoN1E starts n identical processes (i.e. goroutines) within this group
// using a one-argument function returning an optional error.
// This method can be called multiple times with different functions as needed.
// The process argument receives the index in the sequence, starting from one.
func (pg *ProcessGroup) GoN1E(n int, process func(j int) error) {
	pg.goNE(1, n+1, process)
}

func (pg *ProcessGroup) goNE(from, n int, process func(j int) error) {
	pg.joiner.Add(n - from)
	for i := from; i < n; i++ {
		go func(j int) {
			defer pg.joiner.Done()
			if err := process(j); err != nil {
				pg.µ.Lock()
				pg.errs = append(pg.errs, err)
				pg.µ.Unlock()
			}
		}(i)
	}
}

//-------------------------------------------------------------------------------------------------

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

// JoinE is called by the parent goroutine when it wants to sit and wait for
// every process (goroutine) in this group to have terminated. JoinE will therefore
// block until this condition is reached.
//
// See Join for further details.
//
// It returns the collection of all errors that arose from the processes. If JoinE
// is called multiple times, only the newly-arising errors will be returned each
// time.
func (pg *ProcessGroup) JoinE() error {
	pg.Join()

	pg.µ.Lock()
	defer pg.µ.Unlock()

	me := pg.errs.Err()
	pg.errs = nil
	return me
}
