# process - goroutines as contained processes

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/process)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/process)](https://goreportcard.com/report/github.com/rickb777/process)

Simple fork and join of goroutines - easy to use; no fuss.

Puts the 'P' of CSP back into Go.

All it does is **handle the join when a group of goroutines terminate**. Internally, a `sync.WaitGroup` is
administered for you. There's not much to it but it makes your job easier.

You can also limit the maximum number of concurrent goroutines, e.g. to create a worker pool.

## Installation

    go get -u github.com/rickb777/process

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

## WorkQueue

WorkQueue is a function that returns a channel with unlimited buffering. This is useful for work queues
in channel networks that might otherwise deadlock because they contain loops.

A WorkQueue is particularly useful when combined with fixed-size goroutine pools, as described above.

## Hierarchies

A process group contains processes. These processes can also be process groups, or they can contain process
groups. As long as the `Wait` calls are positioned so that each group terminates tidily, the nesting should
*just work* (TM).


## What's New in Version 2

* Several API functions were deleted to keep things simple.
* New MaxConcurrency throttling feature was added.

## Licence

[MIT](LICENSE)
