// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package main

import (
	"os"

	"github.com/nico612/go-project/examples/miniblog/internal/miniblog"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd := miniblog.NewMiniBlogCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
