package main

import (
	"github.com/nico612/go-project/internal/miniblog"
	"os"
)

func main() {
	cmd := miniblog.NewMiniBlogCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}