package process

import (
	"testing"
)

func TestProcessGroup(t *testing.T) {
	n1 := false
	n2 := false
	n3 := false
	n4 := false
	n5 := false

	pg := NewGroup()

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

func TestProcessGroupN(t *testing.T) {
	c := make(chan int)
	pg := NewGroup()

	pg.GoN(2, func() { c <- 1 })

	<-c
	<-c

	pg.Join() // no deadlock expected
}
