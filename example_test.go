package process

import (
	"errors"
	"fmt"
)

func ExampleProcessGroup_GoN() {
	// This example creates five goroutines and receives a number from each one on a channel.
	c := make(chan int)

	pg := NewGroup()

	// the value of i ranges from 0 to 4
	pg.GoN(5, func(i int) { c <- i })

	var sum int
	sum += <-c
	sum += <-c
	sum += <-c
	sum += <-c
	sum += <-c

	pg.Wait() // no deadlock expected

	fmt.Println(pg.Err())
	fmt.Println("sum", sum)
	// Output: <nil>
	// sum 10
}

func ExampleProcessGroup_GoE() {
	// This example creates two goroutines that return errors.
	var (
		e1 = errors.New("E1")
		e2 = errors.New("E2")
	)

	pg := NewGroup()

	pg.GoE(func() error { return e1 })
	pg.GoE(func() error { return e2 })

	pg.Wait()

	combinedErrors := pg.Err()
	// The errors may have arisen in either order so the combined error
	// reflects this non-determinism. Its error message contains both
	// "E1" and "E2".

	fmt.Println(len(combinedErrors.Error()))
	// Output: 5
}

func ExampleProcessGroup_GoNE() {
	// This example creates two goroutines that return errors.
	var (
		e1 = errors.New("E1")
	)

	pg := NewGroup()

	pg.GoNE(5, func(i int) error { return e1 })

	pg.Wait()

	combinedErrors := pg.Err()
	// The errors may have arisen in either order so the combined error
	// reflects this non-determinism. Its error message contains both
	// "E1" and "E2".

	fmt.Println(combinedErrors.Error())
	// Output: E1
	// E1
	// E1
	// E1
	// E1
}
