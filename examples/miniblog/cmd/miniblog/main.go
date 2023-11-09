package main

import (
	"github.com/nico612/go-project/examples/miniblog/internal"
	_ "go.uber.org/automaxprocs"

	"os"
)

func main() {
	cmd := internal.NewMiniBlogCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
