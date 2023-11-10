package main

import (
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog"
	_ "go.uber.org/automaxprocs"

	"os"
)

func main() {
	cmd := miniblog.NewMiniBlogCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
