package process

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

var (
	e1 = errors.New("E1")
	e2 = errors.New("E2")
)

func TestEmpty(t *testing.T) {
	pg := NewGroup()

	// expect no runtime error nor deadlock
	pg.Wait()

	if pg.Err() != nil {
		t.Fatalf("expected no error")
	}
}

func TestProcessGroup_Go_GoE(t *testing.T) {
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		n1, n2, n3, n4, n5 := false, false, false, false, false

		pg.Go(func() { n1 = true })
		pg.Go(func() { n2 = true })
		pg.Go(func() { n3 = true })
		pg.GoE(func() error { n4 = true; return e1 })
		pg.GoE(func() error { n5 = true; return e2 })

		pg.Wait() // no deadlock expected
		e := pg.Err()

		if !n1 {
			t.Errorf("no n1")
		}
		if !n2 {
			t.Errorf("no n2")
		}
		if !n3 {
			t.Errorf("no n3")
		}
		if !n4 {
			t.Errorf("no n5")
		}
		if !n5 {
			t.Errorf("no n5")
		}
		if e == nil {
			t.Fatalf("no err")
		}

		em := e.Error()
		if len(em) < 5 {
			t.Errorf("no count: %q", em)
		}
		if strings.Index(em, "E1") < 0 {
			t.Errorf("no E1: %q", em)
		}
		if strings.Index(em, "E2") < 0 {
			t.Errorf("no E2: %q", em)
		}

		if s := pg.Size(); s != 0 {
			t.Errorf("size is %d", s)
		}
	}
}

// FWIW, the equivalent Occam communication patterns are shown as comments.

// -------------------------------------------------------------------------------------------------
// CHAN INT c:
// INT x:
// SEQ
//
//	PAR
//	  c ! 1
//	  c ! 1
//	c ? x
//	c ? x
func TestProcessGroup_GoN(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.GoN(5, func(i int) { c <- i })

		if s := pg.Size(); s == 0 {
			t.Errorf("size is 0")
		}

		var sum int
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c

		if sum != 10 {
			t.Errorf("Got %d", sum)
		}

		pg.Wait() // no deadlock expected

		if s := pg.Size(); s != 0 {
			t.Errorf("size is %d", s)
		}
	}
}

func TestProcessGroup_GoNE(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.GoNE(5, func(i int) error { c <- i; return e1 })

		var sum int
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c

		if sum != 10 {
			t.Errorf("Got %d", sum)
		}

		pg.Wait() // no deadlock expected

		e := pg.Err()
		if e == nil {
			t.Fatalf("no err")
		}

		em := e.Error()
		if em != "E1\nE1\nE1\nE1\nE1" {
			t.Errorf("bad errors: %q", em)
		}
	}
}

// -------------------------------------------------------------------------------------------------
// CHAN INT c1:
// INT x, y:
// SEQ
//
//	PAR
//	  c1 ! 1
//	  c1 ! 1
//	  SEQ
//	    PAR
//	      c2 ! 1
//	      c2 ! 1
//	    c2 ? y
//	    c2 ? y
//	    c2 ! 1
//	c1 ? x
//	c1 ? x
//	c1 ? x
func TestProcessGroupNested(t *testing.T) {
	c1 := make(chan int)
	c2 := make(chan int)
	pg := NewGroup()
	inner := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.Go(func() {
			c1 <- 1
		})
		pg.Go(func() {
			c1 <- 1
		})
		pg.Go(func() {
			c1 <- 1
			<-c2
			<-c2
			inner.Wait()
		})

		inner.Go(func() {
			c2 <- 1
		})
		inner.Go(func() {
			c2 <- 1
		})

		<-c1
		<-c1
		<-c1

		pg.Wait() // no deadlock expected
	}
}

func TestLongChannel(t *testing.T) {
	in, out := WorkQueue[int](16)
	end := make(chan bool)

	in <- 1
	if <-out != 1 {
		t.Fatalf("expected 1")
	}

	in <- 2
	in <- 3
	in <- 4
	in <- 5

	go func() {
		defer func() { end <- true }()

		var act int
		exp := 2
		for act = range out {
			if act != exp {
				panic(fmt.Sprintf("got %d, expected %d", act, exp))
			}
			exp++
		}
		if act != 999 {
			panic(fmt.Sprintf("ended with %d; expected 999", act))
		}
	}()

	for i := 6; i < 1000; i++ {
		in <- i
	}

	close(in)
	<-end
}
