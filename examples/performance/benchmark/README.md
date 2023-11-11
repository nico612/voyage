
执行命令, 采集性能数据
```shell
go test -bench=BenchmarkExp -benchmem -cpuprofile cpu.profile -memprofile mem.profile
```

```shell
$ go test -bench=BenchmarkExp -benchmem -cpuprofile cpu.profile -memprofile mem.profile 
goos: linux
goarch: amd64
pkg: github.com/marmotedu/miniblogdemo/performance/benchmark
cpu: Intel(R) Xeon(R) Platinum 8255C CPU @ 2.50GHz
BenchmarkExp-16            100000000                11.15 ns/op               0 B/op               0 allocs/op
PASS
ok          github.com/marmotedu/miniblogdemo/performance/benchmark        1.243s
```
命令行参数解析：

- -bench=BenchmarkExp 指定以 benchmark 的方式运行；

- -benchmem 表示在收集 CPU 的信息外还额外收集内存数据；

- -cpuprofile 指定输出 CPU 的 profile；

- -memprofile 指定输出 mem 的 profile。

运行后可以输出单次操作的耗时和内存消耗，而且生成的 profile 文件可以通过 go tool pprof cpu.profile 查看具体信息。

上述命令控制台输出各项指标含义如下：

- BenchmarkExp-16：表示 GOMAXPROC 为 16；

- 100000000：表示共执行 100000000 次；

- 11.15 ns/op：表示每次耗时 11.15 ns；

- 0 B/op：表示每次操作分配 0 字节；

- 0 allocs/op：表示每次操作分配 0 次内存。

上述命令，也是在当前目录会生成 cpu.profile、mem.profile、xxx.test 文件，可通过 `go tool pprof xxx.test xxx.profile `来查看详细信息，输入 `top` 查看 cpu 占比和分配内存占比排名。
```shell
go tool pprof benchmark.test cpu.profile
File: benchmark.test
Type: cpu
Time: Nov 11, 2023 at 3:56pm (CST)
Duration: 1.82s, Total samples = 1.52s (83.68%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1520ms, 100% of 1520ms total
      flat  flat%   sum%        cum   cum%
    1380ms 90.79% 90.79%     1380ms 90.79%  math.archExp
     140ms  9.21%   100%     1520ms   100%  /performance/benchmark.BenchmarkExp
         0     0%   100%     1380ms 90.79%  math.Exp (inline)
         0     0%   100%     1520ms   100%  testing.(*B).launch
         0     0%   100%     1520ms   100%  testing.(*B).runN
(pprof) 

```
可以看出 `math.archExp` 函数占比 90.79%  基准测试函数`benchmark.BenchmarkExp` 占比9.21% , 可以很直观的看出每个函数占用的耗时百分比，从而对占比多的函数进行优化