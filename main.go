package main

import (
	"fmt"
	"os"

	"github.com/kdisneur/<project>/internal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("version is", internal.GetVersionInfo())

	return nil
}
