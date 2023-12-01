/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package signals

import (
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler 注册用于捕获操作系统信号的处理程序，并返回一个通道 stopCh。
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	// 创建一个 stop 通道，当收到注册的信号时，将会关闭这个通道。
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)

	// 注册一系列操作系统信号（比如 SIGTERM 和 SIGINT）到通道 c 中。
	signal.Notify(c, shutdownSignals...)

	// 监听通道 c 上的信号。
	go func() {
		<-c // 当收到第一个信号时，关闭 stop 通道。
		close(stop)
		<-c
		os.Exit(1) // 第二个信号，程序会直接退出，并返回状态码 1。
	}()

	return stop
}
