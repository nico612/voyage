// Package strategy 策略模式
package strategy

// 定义一个策略类

type IStrategy interface {
	do(int, int) int
}

// 策略实现：加
type add struct {
}

func (*add) do(a, b int) int {
	return a + b
}

// 策略实现：减
type reduce struct {
}

func (*reduce) do(a, b int) int {
	return a - b
}

// Operator 具体策略执行者
type Operator struct {
	strategy IStrategy
}

// SetStrategy 设置策略
func (o *Operator) SetStrategy(strategy IStrategy) {
	o.strategy = strategy
}

// 调用策略中的方法
func (o *Operator) calculate(a, b int) int {
	return o.strategy.do(a, b)
}
