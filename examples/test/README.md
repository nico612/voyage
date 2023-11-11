
## 单元测试
参考文章：https://geektutu.com/post/quick-go-test.html

### 规范
- 测试用例名称一般命名为 Test 加上待测试的方法名。
- 测试用的参数有且只有一个，在这里是 t *testing.T。
- 基准测试(benchmark)的参数是 *testing.B，TestMain 的参数是 *testing.M 类型。

### 运行命令

`go test`，该 package 下所有的测试用例都会被执行。

`go test -v `， `-v`：显示每个用例详细测试结果，`-cover`：查看覆盖率

`go test -run TestAdd -v`：只运行 `TestAdd`这一个用例 `-run` 指定运行的用例，参数支持通配符`*`，和部分正则表达式，例如：`^`、`$`


### 子测试（Subtest）
子测试是 Go 语言内置支持的，可以在某个测试用例中，根据测试场景使用 t.Run创建不同的子测试用例

`go test -run TestMul/pos  -v` 指定运行某个测试用例中的子测试


### 帮助函数（helps）
对一些重复的逻辑，抽取出来作为公共的帮助函数(helpers)，可以增加测试代码的可读性和可维护性。 借助帮助函数，可以让测试用例的主逻辑看起来更清晰。

`t.Helper() `

### setup 和 teardown

如果在**同一个测试文件中**，每一个测试用例运行前后的逻辑是相同的，一般会写在 setup 和 teardown 函数中。`setup()`在所有测试执行前调用 `teardown()` 所有测试运行完后调用

例如执行前需要实例化待测试的对象，如果这个对象比较复杂，很适合将这一部分逻辑提取出来；执行后，可能会做一些资源回收类的工作，例如关闭网络连接，释放文件等。


### 网络测试

####  TCP/HTTP
假设需要测试某个 API 接口的 handler 能够正常工作，例如 helloHandler

```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello world"))
}
```
那我们可以创建真实的网络连接进行测试：
```go
// test code
import (
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

func handleError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("failed", err)
	}
}

func TestConn(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	handleError(t, err)
	defer ln.Close()

	http.HandleFunc("/hello", helloHandler)
	go http.Serve(ln, nil)

	resp, err := http.Get("http://" + ln.Addr().String() + "/hello")
	handleError(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handleError(t, err)

	if string(body) != "hello world" {
		t.Fatal("expected hello world, but got", string(body))
	}
}
```

- net.Listen("tcp", "127.0.0.1:0")：监听一个未被占用的端口，并返回 Listener。
- 调用 http.Serve(ln, nil) 启动 http 服务。
- 使用 http.Get 发起一个 Get 请求，检查返回值是否正确。
- 尽量不对 http 和 net 库使用 mock，这样可以覆盖较为真实的场景。

#### httptest

针对 http 开发的场景，使用标准库 net/http/httptest 进行测试更为高效。
```go
// test code
import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConn(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	helloHandler(w, req)
	bytes, _ := ioutil.ReadAll(w.Result().Body)

	if string(bytes) != "hello world" {
		t.Fatal("expected hello world, but got", string(bytes))
	}
}
```
使用 httptest 模拟请求对象(req)和响应对象(w)，达到了相同的目的。

## Benchmark 基准测试
基准测试用例的定义如下：
```go
func BenchmarkName(b *testing.B){
    // ...
}
```
- 函数名必须以 Benchmark 开头，后面一般跟待测试的函数名
- 参数为 b *testing.B。
- 执行基准测试时，需要添加 -bench 参数。

例如：
```go
func BenchmarkHello(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("hello")
    }
}
```
```go
$ go test -benchmem -bench .
...
BenchmarkHello-16   15991854   71.6 ns/op   5 B/op   1 allocs/op
...
```
基准测试报告每一列值对应的含义如下：

```go
type BenchmarkResult struct {
    N         int           // 迭代次数
    T         time.Duration // 基准测试花费的时间
    Bytes     int64         // 一次迭代处理的字节数
    MemAllocs uint64        // 总的分配内存的次数
    MemBytes  uint64        // 总的分配内存的字节数
}
```

如果在运行前基准测试需要一些耗时的配置，则可以使用 b.ResetTimer() 先重置定时器，例如：

```go
func BenchmarkHello(b *testing.B) {
    ... // 耗时操作
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("hello")
    }
}
```

使用 RunParallel 测试并发性能
```go
func BenchmarkParallel(b *testing.B) {
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}
```
```go
$ go test -benchmem -bench .
...
BenchmarkParallel-16   3325430     375 ns/op   272 B/op   8 allocs/op
...

```

## mock/stub 测试，
gomock 是官方提供的 mock 框架，同时还提供了 mockgen 工具用来辅助生成测试代码。

使用如下命令即可安装：

```shell
go get -u github.com/golang/mock/gomock
go get -u github.com/golang/mock/mockgen
```

### 一个简单的 Demo

```go
// db.go
type DB interface {
	Get(key string) (int, error)
}

func GetFromDB(db DB, key string) int {
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}
```

假设 DB 是代码中负责与数据库交互的部分(在这里用 map 模拟)，测试用例中不能创建真实的数据库连接。这个时候，如果我们需要测试 GetFromDB 这个函数内部的逻辑，就需要 mock 接口 DB。

第一步：使用 mockgen 生成 db_mock.go。一般传递三个参数。包含需要被mock的接口得到源文件source，生成的目标文件destination，包名package。
```shell
$ mockgen -source=db.go -destination=db_mock.go -package=main
```
第二步：新建 db_test.go，写测试用例。
```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))

	if v := GetFromDB(m, "Tom"); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
}
```
- 这个测试用例有2个目的，一是使用 ctrl.Finish() 断言 DB.Get() 被是否被调用，如果没有被调用，后续的 mock 就失去了意义；
- 二是测试方法 GetFromDB() 的逻辑是否正确(如果 DB.Get() 返回 error，那么 GetFromDB() 返回 -1)。
- NewMockDB() 的定义在 db_mock.go 中，由 mockgen 自动生成。

执行测试：
```shell
$ go test . -cover -v
=== RUN   TestGetFromDB
--- PASS: TestGetFromDB (0.00s)
PASS
coverage: 81.2% of statements
ok      example 0.008s  coverage: 81.2% of statements

```

### 打桩(stubs)
在上面的例子中，当 Get() 的参数为 Tom，则返回 error，这称之为打桩(stub)，有明确的参数和返回值是最简单打桩方式。除此之外，检测调用次数、调用顺序，动态设置返回值等方式也经常使用。

#### 参数(Eq, Any, Not, Nil)
```go
m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
m.EXPECT().Get(gomock.Any()).Return(630, nil)
m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil) 
m.EXPECT().Get(gomock.Nil()).Return(0, errors.New("nil")) 
```

- Eq(value) 表示与 value 等价的值。
- Any() 可以用来表示任意的入参。
- Not(value) 用来表示非 value 以外的值。
- Nil() 表示 None 值

#### 返回值(Return, DoAndReturn)

```go
m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil)
m.EXPECT().Get(gomock.Any()).Do(func(key string) {
    t.Log(key)
})
m.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (int, error) {
    if key == "Sam" {
        return 630, nil
    }
    return 0, errors.New("not exist")
})
```

#### 调用次数(Times)
```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil).Times(2)
	GetFromDB(m, "ABC")
	GetFromDB(m, "DEF")
}
```
- Times() 断言 Mock 方法被调用的次数。
- MaxTimes() 最大次数。
- MinTimes() 最小次数。
- AnyTimes() 任意次数（包括 0 次）。

#### 调用顺序(InOrder)
```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	o1 := m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
	o2 := m.EXPECT().Get(gomock.Eq("Sam")).Return(630, nil)
	gomock.InOrder(o1, o2)
	GetFromDB(m, "Tom")
	GetFromDB(m, "Sam")
}

```

### 如何编写可 mock 的代码
写可测试的代码与写好测试用例是同等重要的，如何写可 mock 的代码呢？
- mock 作用的是接口，因此将依赖抽象为接口，而不是直接依赖具体的类。
- 不直接依赖的实例，而是使用依赖注入降低耦合性。

```go
func GetFromDB(key string) int {
	db := NewDB()
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}
```
对 DB 接口的 mock 并不能作用于 GetFromDB() 内部，这样写是没办法进行测试的。那如果将接口 db DB 通过参数传递到 GetFromDB()，那么就可以轻而易举地传入 Mock 对象了。

### Mock 工具
Mock 工具用的最多的是 Go 官方提供的 Mock 框架 GoMock。关于 GoMock 的使用方法，可参考：Go Mock (gomock)简明教程。

此外，还有一些其他的优秀 Mock 工具可供我们使用。这些 Mock 工具分别用在不同的 Mock 场景中，常用的 Mock 工具如下：

sqlmock：可以用来模拟数据库连接。数据库是项目中比较常见的依赖，在遇到数据库依赖时都可以用它。

httpmock：可以用来 Mock HTTP 请求。

bouk/monkey：猴子补丁，能够通过替换函数指针的方式来修改任意函数的实现。如果 GoMock、sqlmock 和 httpmock 这几种方法都不能满足我们的需求，我们可以尝试用猴子补丁的方式来 Mock 依赖。可以这么说，猴子补丁提供了单元测试 Mock 依赖的最终解决方案。



参考链接：https://geektutu.com/post/quick-go-test.html

从零开发企业级Go应用：https://juejin.cn/book/7176608782871429175/section/7179878326159278136


