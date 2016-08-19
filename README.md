# process

Simple fork and join of goroutines - easy to use; no fuss.

Puts the 'P' of CSP back into Go.

All it does is handle the join when a group of goroutines terminate.

Internally, a `sync.WaitGroup` is administered for you. There's not much
to it but it makes your job easier.

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
