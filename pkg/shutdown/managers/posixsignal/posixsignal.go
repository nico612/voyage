package posixsignal

import (
	"github.com/nico612/voyage/pkg/shutdown"
	"os"
	"os/signal"
	"syscall"
)

// Name defines shutdown manager name.
const Name = "PosixSignalManager"

// PosixSignalManager 实现了 shutdown.ShutdownManager 接口，它可以被添加到 shutdown.GracefulShutdown 中。
// 使用 NewPosixSignalManager 进行初始化。
type PosixSignalManager struct {
	signals []os.Signal
}

// NewPosixSignalManager 初始化了 PosixSignalManager。
// 可以传入 os.Signal 类型的参数来监听信号
// 如果没有提供参数，则默认监听 SIGINT 和 SIGTERM 信号。
// kill （不带参数）表示发送 syscall.SIGTERM 信号（默认）
// kill -2 表示发送 syscall.SIGINT 信号
// kill -9 表示发送 syscall.SIGKILL 信号, 强制终止进程的信号。它会立即终止进程，而不会给进程任何处理或清理的机会，该信号不会被获取到，所以这里不需要添加该信号
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {

		sig = make([]os.Signal, 2)

		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{signals: sig}
}

func (psm *PosixSignalManager) GetName() string {
	return Name
}

// Start 开始监听信号
func (psm *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, psm.signals...)
		<-c
		gs.StartShutdown(psm)
	}()

	return nil
}

// ShutdownStart does nothing.
func (psm *PosixSignalManager) ShutdownStart() error {
	return nil
}

// ShutdownFinish exits the app with os.Exit(0).
func (psm *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)
	return nil
}
