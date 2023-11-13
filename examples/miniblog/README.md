
# miniblog
教程地址：https://juejin.cn/book/7176608782871429175


## 项目

### 初始化
```shell
make init -f Makefile
```
### build

```shell
make build -f Makefile
```

### 运行
```shell
make run -f Makefile
```

### 打印版本号
```shell

./bin/miniblog/miniblog --version

#打印json格式的版本号
./bin/miniblog/miniblog --version=raw 

```


### 启动HTTPS服务
#### 生成证书
`make ca -f miniblog.mk`

#### 启动
`make run -f miniblog.mk` 

## 开发工具

### Swagger 
OpenAPI 规范（以前称为 Swagger 规范）是 REST API 的 API 描述格式。 OpenAPI 规范兼容文件允许我们描述完整的 REST API。 它通常以 YAML 或 JSON 文件格式编写。目前最新的 OpenAPI 规范是OpenAPI 3.0（也就是 Swagger 2.0 规范）。

#### 安装
```shell
$ go install github.com/go-swagger/go-swagger/cmd/swagger@latest

```
#### 运行
```shell
$ swagger serve -F=swagger --no-open --port 65534 ./api/openapi/openapi.yaml
```
在浏览器中打开 http://localhost:65534/docs


**这里需要注意**：使用 swagger serve 渲染 OpenAPI 文档需要确保 OpenAPI 文档版本为：swagger: "2.0"，例如：
```yaml
swagger: "2.0"
servers:
  - url: http://127.0.0.1:8080/v1
    description: development server
info:
  version: "1.0.0"
  title: miniblog api definition

```
编写后的 OpenAPI 文档需要根据目录规范存放在：api/openapi 目录下。


### Air
程序热加载工具，在 Go 项目开发中，比较受欢迎的是 air 工具。

#### 安装 air 工具。

```shell
go install github.com/cosmtrek/air@latest
```
#### 配置 air 工具。
官方给出的示例配置：[air_example.toml](https://github.com/cosmtrek/air/blob/master/air_example.toml)air_example.toml 里面的示例配置基本能满足绝大部分的项目需求，一般只需要再配置 cmd、bin、args_bin 3 个参数即可。

在 miniblog 项目根目录下创建 .air.toml 文件，

.air.toml 基于 air_example.toml 文件修改了以下参数配置：

```shell
# 只需要写你平常编译使用的 shell 命令。你也可以使用 `make`.
cmd = "make build"
# 由 `cmd` 命令得到的二进制文件名.
bin = "_output/miniblog"
```
参数介绍：
- cmd：指定了监听文件有变化时，air 需要执行的命令，这里指定了 make build 重新构建 miniblog 二进制文件；
- bin：指定了执行完 cmd 命令后，执行的二进制文件，这里指定了编译构建后的二进制文件 _output/miniblog。
- 启动 air 工具。air # 默认使用当前目录下的 .air.toml 配置，你可以通过 `-c` 选项指定配置，例如：`air -c .air.toml`


## db2struct 数据库表自动生成Model

```shell
$ mkdir -p internal/pkg/model
$ cd internal/pkg/model
$ db2struct --gorm --no-json -H 127.0.0.1 -d miniblog -t user --package model --struct UserM -u miniblog -p '12345678' --target=user.go
$ db2struct --gorm --no-json -H 127.0.0.1 -d miniblog -t post --package model --struct PostM -u miniblog -p '12345678' --target=post.go
```
命令行参数：
```shell
$ db2struct  -h
Usage of db2struct:
        db2struct [-H] [-p] [-v] --package pkgName --struct structName --database databaseName --table tableName
Options:
  -H, --host=         Host to check mariadb status of
  --mysql_port=3306   Specify a port to connect to
  -t, --table=        Table to build struct from
  -d, --database=nil  Database to for connection
  -u, --user=user     user to connect to database
  -v, --verbose       Enable verbose output
  --package=          name to set for package
  --struct=           name to set for struct
  --json              Add json annotations (default)
  --no-json           Disable json annotations
  --gorm              Add gorm annotations (tags)
  --guregu            Add guregu null types
  --target=           Save file path
  -p, --password=     Mysql password
  -h, --help          Show usage message
  --version           Show version

```

### mysqldump 导出数据库和表的SQL语句

```shell
$ mysqldump -h127.0.0.1 -uroot --databases miniblog -p'12345678' --add-drop-database --add-drop-table --add-drop-trigger --add-locks --no-data > configs/miniblog.sql
```


## 应用启动构建框架

- pflag：[如何使用Pflag给应用添加命令行标识](https://github.com/marmotedu/geekbang-go/blob/master/%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8Pflag%E7%BB%99%E5%BA%94%E7%94%A8%E6%B7%BB%E5%8A%A0%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A0%87%E8%AF%86.md)
- viper：[配置解析神器-Viper全解；](https://github.com/marmotedu/geekbang-go/blob/master/%E9%85%8D%E7%BD%AE%E8%A7%A3%E6%9E%90%E7%A5%9E%E5%99%A8-Viper%E5%85%A8%E8%A7%A3.md)
- cobra：[现代化的命令行框架-Cobra全解。](https://github.com/marmotedu/geekbang-go/blob/master/%E7%8E%B0%E4%BB%A3%E5%8C%96%E7%9A%84%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A1%86%E6%9E%B6-Cobra%E5%85%A8%E8%A7%A3.md)
  安装： go get -u github.com/spf13/cobra@latest

## 三方库
- 日志库采用： logrus + zap，教程：https://juejin.cn/book/7176608782871429175/section/7176610186029695037

## grpc + grpc-gateway 
[grpc-hello-world：可用的代码（亲测可用）](https://github.com/eddycjy/grpc-hello-world) + [教程（Grpc+Grpc Gateway实践二 有些复杂的Hello World）；](https://segmentfault.com/a/1190000013408485)


## 接口性能测试
https://juejin.cn/book/7176608782871429175/section/7179878428215083063

### Wrk 工具安装
```shell
$ git clone https://github.com/wg/wrk
$ cd wrk/
$ make
$ sudo cp ./wrk /usr/bin
```

### Wrk 使用方法

Wrk 使用起来不复杂，执行 wrk --help 可以看到 wrk 的所有运行参数：

```shell
$ wrk --help
Usage: wrk <options> <url>
  Options:
    -c, --connections <N>  Connections to keep open
    -d, --duration    <T>  Duration of test
    -t, --threads     <N>  Number of threads to use

    -s, --script      <S>  Load Lua script file
    -H, --header      <H>  Add header to request
        --latency          Print latency statistics
        --timeout     <T>  Socket/request timeout
    -v, --version          Print version details

  Numeric arguments may include a SI unit (1k, 1M, 1G)
  Time arguments may include a time unit (2s, 2m, 2h)

```
常用的参数为：

- -t: 线程数（线程数不要太多，是核数的 2 到 4 倍即可，多了反而会因为线程切换过多造成效率降低）；
- -c: 并发数；
- -d: 测试的持续时间，默认为 10s；
- -T: 请求超时时间；
- -H: 指定请求的 HTTP Header，有些 API 需要传入一些 Header，可通过 Wrk 的 -H 参数来传入；
- --latency: 打印响应时间分布；
- -s: 指定 Lua 脚本，Lua 脚本可以实现更复杂的请求。

### Wrk 结果解析

1. 启动 miniblog 服务。
```shell
make build -f Makefile

# 需要将日志输出定位到 /dev/null，打印到标准输出会严重影响接口性能
examples/miniblog/bin/miniblog/miniblog -c examples/miniblog/configs/miniblog.yaml &>/dev/null 
```

2. 打开一个新的 Linux 终端执行 wrk 进行 API 接口压力测试，命令如下：.
```shell
$ wrk -t32 -c1000 -d30s -T30s --latency http://127.0.0.1:8080/healthz
Running 30s test @ http://127.0.0.1:8080/healthz
  32 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    12.49ms   12.47ms 180.14ms   88.44%
    Req/Sec     2.89k   417.30    13.21k    88.66%
  Latency Distribution
     50%   13.14ms
     75%   20.27ms
     90%   26.40ms
     99%   47.91ms
  2766240 requests in 30.10s, 1.20GB read
  Socket errors: connect 3, read 0, write 0, timeout 0
Requests/sec:  91904.09
Transfer/sec:     40.93MB
```
测试输出解读如下。

32 threads and 1000 connections: 用 32 个线程模拟 1000 个连接，分别对应 -t 和 -c 参数。

Thread Stats： 线程统计。

Latency: 响应时间，有平均值、标准偏差、最大值、正负一个标准差占比；
Req/Sec: 每个线程每秒完成的请求数， 同样有平均值、标准偏差、最大值、正负一个标准差占比。
Latency Distribution: 响应时间分布。

50%: 50% 的响应时间为：13.14ms；
75%: 75% 的响应时间为：20.27ms；
90%: 90% 的响应时间为：26.40ms；
99%: 99% 的响应时间为：47.91ms。
2766240 requests in 30.10s, 1.20GB read: 30s 完成的总请求数（2766240）和数据读取量（1.20GB）；

Socket errors: connect 3, read 0, write 0, timeout 0: 错误统计；

Requests/sec: QPS；

Transfer/sec: TPS。



