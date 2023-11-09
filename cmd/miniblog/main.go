package main

import (
	"os"

	_ "go.uber.org/automaxprocs"

	"github.com/nico612/go-project/internal/miniblog"
)

func main() {
	cmd := miniblog.NewMiniBlogCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
