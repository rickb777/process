package process

import "github.com/gammazero/deque"

// WorkQueue creates a channel that behaves as if its buffer were unlimited.
// The input and output ends are returned. The input end should be closed when done;
// this will close the output end.
//
// This long channel is useful as a work queue because items can be added arbitrarily
// to it without deadlock, even if the channel connectivity graph contains loops
// (i.e. contains the risk of deadlock). Each process inserting into the queue will
// never be blocked. However, there is no guarantee that memory won't run out, so
// users must consider avoidance of deadlock in their design.
//
// The channel's buffer is a ring buffer in memory, which is of course not persistent
// across restarts. A minimum capacity can be provided; this can be absent or zero for
// the default size.
func WorkQueue[T any](minimumCapacity ...uint) (chan<- T, <-chan T) {
	in := make(chan T, 1)
	out := make(chan T)

	go func() {
		canIn := in
		var canOut chan T
		var last T
		var buffer deque.Deque[T]
		if len(minimumCapacity) > 0 {
			buffer.SetBaseCap(int(minimumCapacity[0]))
		}

		for canIn != nil || canOut != nil {
			select {
			case v, open := <-canIn:
				if !open {
					canIn = nil // start closing
				} else if canOut == nil {
					last = v
					canOut = out
				} else {
					buffer.PushBack(v)
					canOut = out
				}

			case canOut <- last:
				if buffer.Len() > 0 {
					last = buffer.PopFront()
				} else {
					canOut = nil
				}
			}
		}

		close(out)
	}()

	return in, out
}
