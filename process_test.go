package process

import (
	"errors"
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
	pg.Join()
}

func TestProcessGroup(t *testing.T) {
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		n1, n2, n3, n4, n5 := false, false, false, false, false

		pg.Go(func() { n1 = true })
		pg.Go(func() { n2 = true })
		pg.Go(func() { n3 = true })
		pg.GoE(func() error { n4 = true; return e1 })
		pg.GoE(func() error { n5 = true; return e2 })

		em := pg.JoinE().Error()

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
		if strings.Index(em, "2 errors: ") < 0 {
			t.Errorf("no count: %q", em)
		}
		if strings.Index(em, "E1") < 0 {
			t.Errorf("no E1: %q", em)
		}
		if strings.Index(em, "E2") < 0 {
			t.Errorf("no E2: %q", em)
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
func TestProcessGroupN(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.GoN(5, func() { c <- 1 })

		var sum int
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c

		if sum != 5 {
			t.Errorf("Got %d", sum)
		}

		pg.Join() // no deadlock expected
	}
}

func TestProcessGroupNE(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.GoNE(5, func() error { c <- 1; return e1 })

		var sum int
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c

		if sum != 5 {
			t.Errorf("Got %d", sum)
		}

		em := pg.JoinE().Error() // no deadlock expected

		if em != "5 errors: E1; E1; E1; E1; E1" {
			t.Errorf("bad errors: %q", em)
		}
	}
}

func TestProcessGroupN0(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.GoN0(5, func(i int) { c <- i })

		var sum int
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c

		if sum != 10 {
			t.Errorf("Got %d", sum)
		}

		pg.Join() // no deadlock expected
	}
}

func TestProcessGroupN0E(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.GoN0E(5, func(i int) error { c <- i; return e1 })

		var sum int
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c

		if sum != 10 {
			t.Errorf("Got %d", sum)
		}

		em := pg.JoinE().Error() // no deadlock expected

		if em != "5 errors: E1; E1; E1; E1; E1" {
			t.Errorf("bad errors: %q", em)
		}
	}
}

func TestProcessGroupN1(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.GoN1(5, func(i int) { c <- i })

		var sum int
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c

		if sum != 15 {
			t.Errorf("Got %d", sum)
		}

		pg.Join() // no deadlock expected
	}
}

func TestProcessGroupN1E(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		pg.GoN1E(5, func(i int) error { c <- i; return e1 })

		var sum int
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c
		sum += <-c

		if sum != 15 {
			t.Errorf("Got %d", sum)
		}

		em := pg.JoinE().Error() // no deadlock expected

		if em != "5 errors: E1; E1; E1; E1; E1" {
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
			inner.Join()
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

		pg.Join() // no deadlock expected
	}
}
