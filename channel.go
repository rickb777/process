package process

import "github.com/gammazero/deque"

// WorkQueue creates a channel that behaves as if its buffer were unlimited.
// The input and output ends are returned. The input end should be closed when done;
// this will close the output end.
//
// This long channel is useful as a work queue because items can be added arbitrarily
// to it without deadlock (although there is no guarantee that memory won't run out),
// even if the channel connectivity graph contains loops (i.e. contains the risk of
// deadlock). Each process inserting into the queue will never be blocked.
//
// The channel's buffer is a ring buffer in memory, which is of course not persistent
// across restarts. An initial size can be provided; this can be zero for the default
// size.
func WorkQueue[T any](initialSize uint) (chan<- T, <-chan T) {
	in := make(chan T, 1)
	out := make(chan T)

	go func() {
		canIn := in
		var canOut chan T
		var last T
		buffer := deque.New[T](int(initialSize), int(initialSize))

		for canIn != nil || canOut != nil {
			select {
			case v, open := <-canIn:
				if !open {
					canIn = nil // start closing
				} else if canOut == nil {
					last = v
				} else {
					buffer.PushBack(v)
				}
				canOut = out

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
