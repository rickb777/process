// See https://magefile.org/

//go:build mage

// Build steps for the process API:
package main

import (
	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	if err := sh.RunV("go", "test", "-race", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "test", "-covermode=count", "-coverprofile=process.out", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "tool", "cover", "-func=process.out"); err != nil {
		return err
	}
	if err := sh.RunV("gofmt", "-l", "-w", "-s", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return err
	}
	return nil
}
