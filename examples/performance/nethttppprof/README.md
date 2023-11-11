
## 通过 net/http/pprof 包采集
如果你想测试 Web 服务的接口性能，通常可以通过在 Go 代码中导入 net/http/pprof 包，之后可以通过浏览器、wget、go tool pprof 等方式来动态收集性能数据。

执行以下命令运行程序：
```shell
 go run main.go
```
之后，可以打开一个新的终端，执行以下命令，来动态收集性能数据：

- 采集 CPU 数据： go tool pprof ``http://localhost:6060/debug/pprof/profile；

- 采集正在使用的内存数据： go tool pprof ``http://localhost:6060/debug/pprof/heap 采集正在使用的内存数据；

- 采集累计使用的内存数据： go tool pprof -alloc_space ``http://localhost:6060/debug/pprof/heap；

- 采集协程数据： go tool pprof ``http://localhost:6060/debug/pprof/goroutine；

- 采集 trace 数据：wget ``http://localhost:6060/debug/pprof/trace?seconds=5`` -O nethttppprof.trace。

之后，就可以使用 pprof 工具进行性能数据查看并分析。


## 数据分析方式

- 有了性能数据，接下来，我们就可以进行性能分析，并进行性能优化。有以下 3 种方式可以用来进行性能分析：

- 使用所采集的 profile 进行分析；

- 使用所采集的 trace 进行分析；

- 通过火焰图进行分析。

### 使用所采集的 profile 进行分析

使用 go tool pprof xxx.profile 命令打开 profile，然后可以在 pprof 命令的交互界面进行性能分析：

- svg：生成函数调用的消耗图（需要安装 graphviz），比如 CPU 占比、分配内存占比、协程占比等（与采集的数据有关）;
- traces：可以打印出调用堆栈中的资源消耗情况；
- list：用来搜索相关模块的每行代码消耗资源情况；
- top：查看消耗最高的地方；
- tree：包含各个主要函数的消耗情况以及子函数的消耗。
  通过 svg 图，可以直观看出哪些地方占比高（比如 mallocgc、syscall 等），这些就是后续需要重点关注的点。

也可在 `go tool pprof -http=0.0.0.0:6060 xxxx.profile` 来启动一个 Web 服务，通过浏览器访问分析结果。

### 使用所采集的 trace 进行分析

执行 `go tool trace -http=0.0.0.0:6061 xxx.trace`，然后通过浏览器打开，之后进行操作。可以通过 trace 查看协程情况、调度情况、具体执行细节等。可用来分析执行延迟和并发问题。

trace 分析方法可参考：[go的请求追踪神器go tool trace。](https://zhuanlan.zhihu.com/p/377145725)

基于通过 **net/http/pprof** 包采集一节中的测试代码，启动一个 Web 服务，并在一个新的 Linux 终端中，获取 trace 数据，并打开浏览器，查看 trace 的详细信息，操作如下：

1. 启动web服务
   `go run main.go`
2. 获取并查看trace信息
   ```shell
    wget http://localhost:6060/debug/pprof/trace?seconds=5 -O nethttppprof.trace
    go tool trace -http=0.0.0.0:6061 nethttppprof.trace
    ```

然后通过浏览器打开 http://xx.xx.xx.xx:6061 来查看 trace 的详细信息。

### 通过火焰图进行分析

我们可以使用 go-torch 工具更直观地查看性能数据。go-torch 是 Uber 公司开源的一款针对 Go 程序的火焰图生成工具，能收集 stack traces，并把它们整理成火焰图，直观地呈现给开发人员。go-torch 是基于使用 BrendanGregg 创建的火焰图工具生成直观的图像，很方便地分析 Go 的各个方法所占用的 CPU 的时间。

go-torch 会通过 pprof 生成火焰矢量图 torch.svg，然后使用浏览器打开查看。操作命令如下：
```shell
$ git clone https://github.com/brendangregg/FlameGraph.git
$ sudo cp  FlameGraph/flamegraph.pl /usr/local/bin
$ flamegraph.pl -h
$ go install github.com/uber/go-torch@latest
$ go-torch -h
$ go-torch -u http://localhost:6060 -t 30 # 确保 `通过 net/http/pprof 包采集` 步骤中的 Web 服务正在运行
INFO[15:17:54] Run pprof command: go tool pprof -raw -seconds 30 http://localhost:6060/debug/pprof/profile
INFO[15:18:24] Writing svg to torch.svg
```
go-torch 会通过 pprof 生成火焰矢量图 torch.svg，然后使用浏览器打开查看




