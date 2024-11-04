// Package process wraps goroutines with the necessary synchronisation so that the caller
// can easily wait on completion.
//
// Following the CSP practice, these goroutines are called 'processes' - the 'P' in
// Communicating Sequential Processes. But this term is not to be confused with other
// usages, especially of processes within operating systems; the latter are unrelated.
package process

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ProcessGroup allows a related group of processes (i.e. goroutines) - zero or
// more - to be launched. The parent goroutine can then wait for completion of the
// entire group via `Wait()`.
//
// A single parent goroutine will own each ProcessGroup. They should not be shared
// by more than one parent.
type ProcessGroup struct {
	joiner sync.WaitGroup
	size   atomic.Int64
	errs   *list[error]
}

// NewGroup creates a new empty process group. Use Go, GoN, GoE and GoNE to start processes
// (i.e. goroutines) within the group.
//
// The maximum number of goroutines can be capped using MaxConcurrency.
func NewGroup() *ProcessGroup {
	pg := &ProcessGroup{
		errs: newList[error](),
	}
	return pg
}

//-------------------------------------------------------------------------------------------------

// Go starts a single process (i.e. goroutine) within this group
// using a zero-argument function.
// This method can be called multiple times with different functions as needed.
//
// Use Wait to wait for all the processes to terminate.
func (pg *ProcessGroup) Go(process func()) {
	pg.GoNE(1,
		func(_ int) error {
			process()
			return nil
		})
}

//-------------------------------------------------------------------------------------------------

// GoN starts n identical processes (i.e. goroutines) within this group
// using a one-argument function.
// This method can be called multiple times with different functions as needed.
// The process argument receives the index from 0 to n-1.
//
// Use Wait to wait for all the processes to terminate.
func (pg *ProcessGroup) GoN(n int, process func(int)) {
	pg.GoNE(n,
		func(j int) error {
			process(j)
			return nil
		})
}

//-------------------------------------------------------------------------------------------------

// GoE starts a single process (i.e. goroutine) within this group
// using a zero-argument function that can return an error.
// This method can be called multiple times with different functions as needed.
//
// Use Wait to wait for all the processes to terminate.
func (pg *ProcessGroup) GoE(process func() error) {
	pg.GoNE(1,
		func(_ int) error {
			return process()
		})
}

//-------------------------------------------------------------------------------------------------

// GoNE starts n identical processes (i.e. goroutines) within this group
// using a one-argument function that can return an error.
// This method can be called multiple times with different functions as needed.
// The process argument receives the index from 0 to n-1.
//
// Use Wait to wait for all the processes to terminate.
func (pg *ProcessGroup) GoNE(n int, process func(j int) error) {
	pg.joiner.Add(n)
	pg.size.Add(int64(n))

	for i := 0; i < n; i++ {
		go func(j int) {
			defer pg.joiner.Done()
			defer pg.size.Add(-1)

			if err := process(j); err != nil {
				pg.errs.Add(err)
			}
		}(i)
	}
}

//-------------------------------------------------------------------------------------------------

// Wait is called by the parent goroutine when it wants to sit and wait for
// every process (goroutine) in this group to have terminated. Wait will therefore
// block until this condition is reached.
//
// Because the process group does not control the internal behaviour of each child
// process (goroutine), it has no means to guarantee that they will all terminate.
// So it is possible for this method to wait forever (deadlock), as a program error.
// It is up to the client code to prevent this by ensuring that all the child
// processes (goroutines) terminate cleanly.
func (pg *ProcessGroup) Wait() {
	pg.joiner.Wait()
}

//-------------------------------------------------------------------------------------------------

// Size returns the current number of processes (i.e. goroutines) that have been created and not
// yet finished. This is an instantaneous value that may change frequently. The value will be zero
// by the time Wait returns.
func (pg *ProcessGroup) Size() int {
	return int(pg.size.Load())
}

//-------------------------------------------------------------------------------------------------

// Err returns all errors that arose from the processes, combined into a single error.
//
// Each time it is called, the collection of errors is cleared, so if Err
// is called multiple times, only the newly-arising errors will be returned each
// time.
//
// Ths simplest use-case is to call this after Wait() just once.
func (pg *ProcessGroup) Err() error {
	return errors.Join(pg.errs.Clear()...)
}
