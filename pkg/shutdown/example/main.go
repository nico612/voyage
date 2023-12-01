package main

import (
	"fmt"
	"github.com/nico612/voyage/pkg/shutdown"
	"github.com/nico612/voyage/pkg/shutdown/managers/posixsignal"
	"time"
)

func main() {

	// 初始化一个 shutdown
	gs := shutdown.New()

	// 添加 shutdown manager
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// 错误处理回调
	gs.SetErrorHandler(shutdown.ErrorFunc(func(err error) {
		fmt.Println("Error:", err)
	}))

	// 添加关闭回调
	gs.AddShutdownCallback(shutdown.ShutdownFunc(func(shutdownManager string) error {
		fmt.Println("Shutdown callback start")
		if shutdownManager == posixsignal.Name {
			// do something
		}
		fmt.Println("Shutdown callback finished")
		return nil
	}))

	// 开始监听
	if err := gs.Start(); err != nil {
		fmt.Println("Start:", err)
		return
	}

	// do other stuff
	time.Sleep(time.Hour * 2)

}
