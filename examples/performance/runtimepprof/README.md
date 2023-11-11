
## runtime/pprof 包采集

通过在 Go 代码中使用 runtime/pprof 包，可以采集 CPU 信息、内存信息、协程信息等。生成的 profile 文件可以使用 go tool pprof xxx去查看详细信息。

执行以下命令收集性能数据：
```shell
$ go run main.go -cpuprofile=cpu.profile -memprofile=mem.profile
```
以上命令会在当前目录下生成 cpu.profile 和 mem.profile 文件。之后，你可以使用 go tool pprof xxx去查看对应的性能数据。

`go tool pprof mem.profile cpu.profile`