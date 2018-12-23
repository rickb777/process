package process

import (
	"testing"
)

func TestProcessGroup(t *testing.T) {
	pg := NewGroup()

	// expect to be able to make more than one pass
	for p := 0; p < 100; p++ {
		n1, n2, n3, n4, n5 := false, false, false, false, false

		pg.Go(func() { n1 = true })
		pg.Go(func() { n2 = true })
		pg.Go(func() { n3 = true })
		pg.Go(func() { n4 = true })
		pg.Go(func() { n5 = true })

		pg.Join()

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
	}
}

// FWIW, the equivalent Occam communication patterns are shown as comments.

//-------------------------------------------------------------------------------------------------
// CHAN INT c:
// INT x:
// SEQ
//   PAR
//     c ! 1
//     c ! 1
//   c ? x
//   c ? x
func TestProcessGroupN0(t *testing.T) {
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

//-------------------------------------------------------------------------------------------------
// CHAN INT c1:
// INT x, y:
// SEQ
//   PAR
//     c1 ! 1
//     c1 ! 1
//     SEQ
//       PAR
//         c2 ! 1
//         c2 ! 1
//       c2 ? y
//       c2 ? y
//       c2 ! 1
//   c1 ? x
//   c1 ? x
//   c1 ? x
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
