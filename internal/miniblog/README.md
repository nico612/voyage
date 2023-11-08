
# miniblog
教程地址：https://juejin.cn/book/7176608782871429175


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



## 应用启动构建框架

- pflag：[如何使用Pflag给应用添加命令行标识](https://github.com/marmotedu/geekbang-go/blob/master/%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8Pflag%E7%BB%99%E5%BA%94%E7%94%A8%E6%B7%BB%E5%8A%A0%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A0%87%E8%AF%86.md)
- viper：[配置解析神器-Viper全解；](https://github.com/marmotedu/geekbang-go/blob/master/%E9%85%8D%E7%BD%AE%E8%A7%A3%E6%9E%90%E7%A5%9E%E5%99%A8-Viper%E5%85%A8%E8%A7%A3.md)
- cobra：[现代化的命令行框架-Cobra全解。](https://github.com/marmotedu/geekbang-go/blob/master/%E7%8E%B0%E4%BB%A3%E5%8C%96%E7%9A%84%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A1%86%E6%9E%B6-Cobra%E5%85%A8%E8%A7%A3.md)
  安装： go get -u github.com/spf13/cobra@latest

