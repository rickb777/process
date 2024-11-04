# process - goroutines as contained processes

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/process/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/process/v2)](https://goreportcard.com/report/github.com/rickb777/process/v2)

Simple fork and join of goroutines - easy to use; no fuss.

Puts the 'P' of CSP back into Go.

All it does is **handle the join when a group of goroutines terminate**. Internally, a `sync.WaitGroup` is
administered for you. There's not much to it but it makes your job easier.

You can also limit the maximum number of concurrent goroutines, e.g. to create a worker pool.

## Installation

    go get -u github.com/rickb777/process/v2

## ProcessGroup

Just create a new group then tell some functions to `Go`:

```
	processes := process.NewGroup()
	processes.Go(func() {
		...  some work, just a normal goroutine function
	})
	processes.Go(func() {
		...  some other work
	})
	processes.Wait()
```

Another useful case is to create a fixed-size pool of goroutines. How to do this is shown in the
other examples in the documentation. This is an easy and simple way to limit concurrency for some
reason.

The process function can be a closure referring to variables in scope. Channels are a useful way
to return results.

## WorkQueue - channel with a very long buffer

WorkQueue is a function that returns a channel with unlimited buffering. This is useful for work queues
in channel networks that might otherwise deadlock because they contain loops.

A WorkQueue is particularly useful when combined with fixed-size goroutine pools, as described above.

## Hierarchies

A process group contains processes. These processes can also be process groups, or they can contain process
groups. As long as the `Wait` calls are positioned so that each group terminates tidily, the nesting will
*just work*.

## What's New in Version 2

* Several API functions were deleted to keep things simple.
* New WorkQueue unlimited channel was added.

## Licence

[MIT](LICENSE)
