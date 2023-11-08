// Package option 选项模式
package option

import "time"

const (
	defaultTimeout = 10
	defaultCaching = false
)

type Connection struct {
	addr    string
	cache   bool
	timeout time.Duration
}

type options struct {
	caching bool
	timeout time.Duration
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithTimeout(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.timeout = t // 自定义赋值逻辑
	})
}

func WithCaching(caching bool) Option {
	return optionFunc(func(o *options) {
		o.caching = caching // 自定义赋值逻辑
	})
}

func NewConnect(addr string, opts ...Option) (*Connection, error) {

	// 优先创建一个带有默认值的options
	options := options{
		caching: defaultCaching,
		timeout: defaultTimeout,
	}

	for _, opt := range opts {
		opt.apply(&options) // 通过该方法来修改默认值的变量
	}

	return &Connection{
		addr:    addr,
		cache:   options.caching,
		timeout: options.timeout,
	}, nil
}
