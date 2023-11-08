设计模式
设计模式可分为：创建型模式、行为行模式、结构型模式

- 创建型模式：掌握单例模式、简单工厂模式、抽象工厂模式、工厂方法模式
- 行为型模式：掌握策略模式、模版模式
- 结构型模式：掌握代理模式、选项模式

## 创建型模式
提供一种创建对象的同时隐藏创建逻辑的方式，而不是使用new运算符直接实例化对象

### 单例模式

全局只有一个实例，并且它负责创建自己的对象。全局共享一个实例、且只需要被初始化一次。

优点：减少内存和系统的资源的开销、防止多个实例产生冲突等

使用场景：数据库实例、全局配置、全局任务池等。

#### 饿汉式
全局的单例实例在包被加载时创建

```go
package singleton
// 饿汉式单例模式

type singleton struct {
}

var ins *singleton = &singleton{}

func GetInsOr() *singleton {
	return ins
}

```
**注意：** 因为实例是在包被导入时初始化的，所以如果初始化时间增加，会导致程序加载时间比较长

#### 懒汉式
全局单例在第一次被使用时创建，也是使用最多的单例模式
```go
package singleton

import "sync"

type lazySingleton struct {
}

var lazyIns *lazySingleton
var once sync.Once

func GetLazyIns() *lazySingleton {
	
	once.Do(func() { // 使用once.Do 确保实例全局只被创建一次，还可确保当同时叉棍见多个动作时，只有一个动作被执行
		lazyIns = &lazySingleton{}
	})
	return lazyIns
}
```
### 工厂模式
```go
package factory

import "fmt"

type Person struct {
	Name string
	Age int
}

func (p Person) Greet()  {
	fmt.Printf("Hi! My name is %s,", p.Name)
}
```
#### 简单工厂模式
```go
package factory

// NewPerson 简单工厂模式
func NewPerson(name string, age int) *Person {
	return &Person{
		Name: name,
		Age: age,
	}
}

```
简单工厂模式：确保创建的实例具有需要的参数，进而保证实例的方法可以按预期执行

#### 抽象工厂模式
和简单工厂模式的区别：返回的是接口而不是结构体
```go
package abstract

import "fmt"

type Person interface {
	Greet()
}

type person struct {
	name string
	age  int
}

func (p person) Greet() {
	fmt.Printf("Hi! My name is %s,", p.name)
}

// NewPerson 抽象工厂模式：返回一个接口
func NewPerson(name string, age int) Person {
	return person{
		name: name,
		age:  age,
	}
}

```
通过返回接口，还可以使用多个工厂函数来返回不同接口的实现，例如：
```go
package abstract

import (
	"net/http"
	"net/http/httptest"
)

// Doer 定义一个Doer接口，该接口具有一个Do方法
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewHTTPClient 返回一个net/http 包提供的HTTP客户端
func NewHTTPClient() Doer {
	return &http.Client{}
}

// mock客户端
type mockHTTPClient struct{}

// NewMockHTTPClient 返回一个模拟的HTTP客户端，
func NewMockHTTPClient() Doer {
	return &mockHTTPClient{}
}

// Do 该HTTP客户端接收任何请求，并返回一个空的响应
func (*mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	res := httptest.NewRecorder()
	return res.Result(), nil
}

// QueryUser 测试案例, 假如要测试下面这段代码
func QueryUser(doer Doer) error {
	req, err := http.NewRequest("Get", "http://baidu.com", nil)
	if err != nil {
		return err
	}

	_, err = doer.Do(req)
	if err != nil {
		return err
	}

	return nil
}

```
测试用例：
```go
package abstract

import "testing"

// 测试用例
func TestQueryUser(t *testing.T) {
	doer := NewMockHTTPClient()
	if err := QueryUser(doer); err != nil {
		t.Errorf("QueryUser failed, err : %v", err)
	}
}

```

#### 工厂方法模式
在简单工厂模式中，依赖于唯一的工厂对象，如果需要实例化一个产品，就要向工厂中传入一个参数，获取对应的对象。如果要增加一种产品，就要在工厂中修改创建产品的函数，这回导致耦合性过高。此时可以使用**工厂方法模式**
```go
package method

type Person struct {
	name string
	age int
}

func NewPersonFactory(age int) func(name string) Person {
	return func(name string) Person {
		return Person{
			name: name,
			age: age,
		}
	}
}

```
然后就可以使用此功能来创建具有默认年龄的工厂：
```go

newBaby := NewPersonFactory(1)
baby := newBaby("john")

newTeenager := NewPersonFactory(16)
teen := newTeenager("jill)

```

## 行为模式

### 策略模式
策略模式定义一组算法，将每个算法都封装起来，并且使他们之间可以互换

在开发中，经常要根据不同的场景，采取不同的措施，也就是不同策略。假设我们需要对a、b这两个整数进行计算，根据条件不同，需要执行不同的计算方式，我们可以把所有的操作都封装在同一个函数中，然后通过if ... else ...的形式来调用不同的计算方式，这种方式称为：**硬编码**。

在实际开发中，随着功能和体验的不断增长，我们需要经常调价/修改策略，进而需要不断修改已有代码，这不仅会让这个函数越来越难维护，还会因为修改带来一些bug。因此为了解耦，我们需要使用策略模式，定义一些独立的类来封装不同的算法，每一个类封装一个具体的算法（即策略）

```go
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

```

测试用例
```go
package strategy

import (
	"fmt"
	"testing"
)

func TestOperator_calculate(t *testing.T) {

	o := &Operator{}
	o.SetStrategy(&add{})
	result := o.calculate(1, 2)
	fmt.Println("add: ", result)

	o.SetStrategy(&reduce{})
	result = o.calculate(2, 1)
	fmt.Println("reduce: ", result)

}

```

### 模版模式
定义一个操作中的算法的骨架，并将一些步骤延迟到子类中，这种方法可以让子类不改变一个算法结构的情况下，重新定义该算法的某些步骤

简单来说：模版模式就是将一个类中能够公共使用的方法放置在抽象类中实现，将不能公共使用的方法作为抽象方法，强制子类去实现，这样就能做到了将一个类作为一个模版，让开发者去填充需要填充的地方。

```go
package template

import "fmt"

type Cooker interface {
	fire()
	cooke()
	outfire()
}

// 类似于一个抽象类
type CookMenu struct {
}

func (CookMenu) fire() {
	fmt.Println("开火")
}

// 做菜，交给具体的子类实现
func (CookMenu) cooke() {

}

func (CookMenu) outfire() {
	fmt.Println("关火")
}

// 封装具体步骤
func doCook(cook Cooker) {
	cook.fire()
	cook.cooke()
	cook.outfire()
}

type XiHongShi struct {
	CookMenu
}

// 子类具体实现,特定方法
func (*XiHongShi) cooke() {
	fmt.Println("做西红柿")
}

type ChaoJiDan struct {
	CookMenu
}

// 子类具体实现,特定方法
func (*ChaoJiDan) cooke() {
	fmt.Println("炒鸡蛋")
}

```

测试用例：

```go
package template

import (
	"fmt"
	"testing"
)

func Test_doCook(t *testing.T) {
	xihongshi := &XiHongShi{}
	doCook(xihongshi)

	fmt.Println("\n 做另外一道菜")

	chaojidan := &ChaoJiDan{}
	doCook(chaojidan)
}

```

## 结构型模式
特点：关注对象之间的通信

### 代理模式 Proxy Pattern
代理模式：可以为另外一个对象提供一个替身或者占位符，以控制对这个对象的访问。
```go
// Package proxy 代理模式
package proxy

import "fmt"

type Seller interface {
	sell(name string)
}

// Station 火车站
type Station struct {
	stock int // 库存
}

func (s *Station) sell(name string) {
	if s.stock > 0 {
		s.stock --
		fmt.Printf("代理点中：%s 买了一张票，剩余：%d \n", name, s.stock)
	} else {
		fmt.Println("票已卖空")
	}
}

// StationProxy 火车站代理点
type StationProxy struct {
	station *Station // 持有一个火车站对象
}

func (proxy *StationProxy) sell(name string) {
	if proxy.station.stock > 0 {
		proxy.station.stock --
		fmt.Printf("代理点中：%s 买了一张票，剩余：%d \n", name, proxy.station.stock)
	} else {
		fmt.Println("票已卖空")
	}
}


```

### 选项模式 Option Pattern
选项模式也是go语言中常用的模式，使用选项模式，我们可以创建一个带有默认值的struct变量，并选择性地修改其中一些参数的值。

go语言中因为不支持给参数设置默认值，为了既能够创建带默认值的实例，又能够创建自定义参数的实例，需要创建一个带默认值的选项，并用该选项创建实例。

```go
// Package option 选项模式
package option

import "time"

const (
	defaultTimeout = 10
	defaultCaching = false
)

type Connection struct {
	addr string
	cache bool
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
		o.timeout = t	// 自定义赋值逻辑
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

```

适用场景：
- 结构体参数很多，创建结构体时，期望创建一个携带默认值的结构体变量，并选择性修改其中一些参数的值
- 结构体参数经常变动，但我们又不想在参数变动时修改创建实例的函数，例如：在结构体中新增一个retry参数，但是不想在NewConnect入参列表中添加retry int这样的参数声明
如果结构体参数比较少，可以慎重考虑要不要采用选项模式
