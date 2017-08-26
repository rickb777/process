# process - goroutines as contained processes

[![GoDoc](https://godoc.org/github.com/rickb777/process?status.svg)](https://godoc.org/github.com/rickb777/process)
[![Build Status](https://api.travis-ci.org/rickb777/process.png)](https://travis-ci.org/rickb777/process/builds)
[![Coverage Status](https://coveralls.io/repos/github/rickb777/process/badge.svg?branch=master&service=github)](https://coveralls.io/github/rickb777/process?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/process)](https://goreportcard.com/report/github.com/rickb777/process)

Simple fork and join of goroutines - easy to use; no fuss.

Puts the 'P' of CSP back into Go.

All it does is **handle the join when a group of goroutines terminate**. Internally, a `sync.WaitGroup` is
administered for you. There's not much to it but it makes your job easier.

## Installation

    go get -u github.com/rickb777/process

## How To

Just create a new group then tell some functions to Go:

```
	processes := process.NewGroup()
	processes.Go(func() {
		...  some work, just a normal goroutine function
	})
	processes.Go(func() {
		...  some other work
	})
	processes.Join()
```

Or mix this with a pool of several identical goroutines using GoN:

```
	processes := process.NewGroup()
	processes.GoN(3, func() {
		...  some work, just a normal goroutine function
	})
	processes.Join()
```

That's it.

## Licence

[MIT](LICENSE)
