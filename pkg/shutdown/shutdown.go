package shutdown

import "sync"

// ShutdownCallback 接口用于实现在请求关闭时的回调函数。OnShutdown 方法会在关闭管理器请求关闭时调用，传入的参数是发出关闭请求的 ShutdownManager 的名称。
type ShutdownCallback interface {
	OnShutdown(string) error
}

// ShutdownFunc 一个辅助类型，可以使用匿名函数作为 ShutdownCallbacks。
type ShutdownFunc func(string) error

// OnShutdown 定义了在触发关闭时需要运行的操作。
func (f ShutdownFunc) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

// ShutdownManager 定义关闭管理这所需要的功能
type ShutdownManager interface {

	// GetName 返回 ShutdownManager 的名称。
	GetName() string

	// Start 启动 ShutdownManager 并开始监听关闭请求。
	Start(gs GSInterface) error

	// ShutdownStart 在开始关闭时调用，执行关闭前的操作。
	ShutdownStart() error

	// ShutdownFinish 在关闭完成时调用，执行关闭后的操作。
	ShutdownFinish() error
}

// ErrorHandler 处理异步错误的接口。
// 当异步操作发生错误时，可以实现这个接口，然后将其传递给 SetErrorHandler 方法，以便在出现错误时进行相应的处理。
type ErrorHandler interface {
	OnError(err error)
}

// ErrorFunc 一个辅助类型，它允许您将匿名函数作为 ErrorHandlers 使用。这样，您可以更方便地定义在出现错误时要执行的操作。
type ErrorFunc func(err error)

// OnError 定义在发生错误时执行的操作
func (f ErrorFunc) OnError(err error) {
	f(err)
}

// GSInterface 是由 GracefulShutdown 实现的接口，它被传递给 ShutdownManager，在请求关闭时调用 StartShutdown。
type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallback ShutdownCallback)
}

// GracefulShutdown 是处理 ShutdownCallbacks 和 ShutdownManagers 的主要结构体。使用 New 进行初始化。
type GracefulShutdown struct {
	callbacks    []ShutdownCallback
	managers     []ShutdownManager
	errorHandler ErrorHandler
}

func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]ShutdownCallback, 0, 10),
		managers:  make([]ShutdownManager, 0, 3),
	}
}

// Start 方法调用了所有添加的 ShutdownManager 的 Start 方法。这些 ShutdownManager 开始监听关闭请求。如果任何 ShutdownManager 返回错误，则返回错误。
func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}
	return nil
}

// AddShutdownManager 添加一个 ShutdownManager，这个 ShutdownManager 将会监听关闭请求。
func (gs *GracefulShutdown) AddShutdownManager(manager ShutdownManager) {
	gs.managers = append(gs.managers, manager)
}

// AddShutdownCallback 方法用于添加一个 ShutdownCallback，当关闭请求被触发时将会调用它。
// 你可以提供任何实现了 ShutdownCallback 接口的内容，或者可以提供一个函数，如下所示：
//
//	AddShutdownCallback(shutdown.ShutdownFunc(func() error {
//	   // callback code
//	   return nil
//	}))
func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

// SetErrorHandler 方法用于设置一个 ErrorHandler，当在 ShutdownCallback 或 ShutdownManager 中遇到错误时将会调用它。
// 你可以提供任何实现了 ErrorHandler 接口的内容，或者可以提供一个函数，如下所示：
//
//	SetErrorHandler(shutdown.ErrorFunc(func (err error) {
//		// handle error
//	}))
func (gs *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gs.errorHandler = errorHandler
}

// StartShutdown 方法被 ShutdownManager 调用，并将启动关闭流程。首先调用 ShutdownManager 的 ShutdownStart 方法，然后调用所有 ShutdownCallbacks，等待回调完成，并在 ShutdownManager 上调用 ShutdownFinish 方法。
func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {

	gs.ReportError(sm.ShutdownStart())

	var wg sync.WaitGroup
	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallbak ShutdownCallback) {
			defer wg.Done()
			gs.ReportError(shutdownCallback.OnShutdown(sm.GetName()))
		}(shutdownCallback)
	}

}

// ReportError 用于将错误报告给 ErrorHandler。它在 ShutdownManagers 中使用。如果发生错误且 ErrorHandler 不为 nil，则会将错误传递给 ErrorHandler 处理。
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}
