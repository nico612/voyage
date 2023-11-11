GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)

# VERSION:=$(shell git describe --tags --always)
## 指定应用使用的 version 包，会通过 `-ldflags -X` 向该包中指定的变量注入值
VERSION_PACKAGE=github.com/nico612/go-project/pkg/version

## 定义 VERSION 语义化版本号
ifeq ($(origin VERSION), undefined)
### 获取版本号。
### --tags ：使用所有的标签，而不是只使用带注释的标签（annotated tag）。git tag <tagname> 生成一个不带注释的标签，git tag -a <tagname> -m '<message>'生成一个带注释的标签；
### --always：如果仓库没有可用的标签，那么使用 commit 缩写来替代标签；
### --match <pattern>：只考虑与给定模式相匹配的标签。
VERSION := $(shell git describe --tags --always --match='v*')
endif

## 检查代码仓库是否是 dirty（默认dirty）
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
  GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD) ### git rev-parse HEAD 获取构建时的 Commit ID；

GO_LDFLAGS += \
  -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
  -X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
  -X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
  -X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') ### date -u +'%Y-%m-%dT%H:%M:%SZ' 命令获取构建时间；


# ./ 根目录
# MAKEFILE_LIST 变量是 Makefile 的内置变量，表示：make 所需要处理的 makefile 文件列表，当前 makefile 的文件名总是位于列表的最后，文件名之间以空格进行分隔；
# 函数 $(lastword <text>) 取字符串 <text> 中的最后一个单词，并返回字符串 <text> 的最后一个单词；
# 函数 $(dir <names...>) 从文件名序列 <names> 中取出目录部分。目录部分是指最后一个反斜杠（/）之前的部分。如果没有反斜杠，那么返回 ./；
COMMON_SELF_DIR:=$(dir $(lastword $(MAKEFILE_LIST)))
# 获取项目根目录绝对路径。
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))/examples/miniblog

# bin 目录
BIN_DIR := $(ROOT_DIR)/bin/miniblog

# 构建产物，临时文件存放目录
OUTPUT_DIR := $(ROOT_DIR)/_output

# Protobuf 文件存放路径
APIROOT=$(ROOT_DIR)/pkg/proto

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git | grep cmd))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: test
.test:
	echo $(ROOT_DIR)

.PHONY: all
all: format build # 指定执行 make 命令时默认需要执行的规则目标

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.15.1

.PHONY: build
build: tidy
	mkdir -p bin
	# go build -ldflags "-X main.Version=$(VERSION)" -o $(BIN_DIR)/miniblog $(ROOT_DIR)/cmd/miniblog/main.go
	go build -v -ldflags "$(GO_LDFLAGS)" -o $(BIN_DIR)/miniblog $(ROOT_DIR)/cmd/miniblog/main.go

.PHONY: run
run: build
	$(BIN_DIR)/miniblog -c $(ROOT_DIR)/configs/miniblog.yaml

.PHONY: generate
# generate
generate: tidy
	go generate $(ROOT_DIR)/...

.PHONY: clean # 实现幂等删除。
clean:
	-rm -vrf bin
	-rm -vrf api/doc
	-rm -vrf api/gen

.PHONY: ca
ca: ## 生成 CA 文件
	@mkdir -p $(OUTPUT_DIR)/cert
	@openssl genrsa -out $(OUTPUT_DIR)/cert/ca.key 1024 # 生成根证书私钥
	@openssl req -new -key $(OUTPUT_DIR)/cert/ca.key -out $(OUTPUT_DIR)/cert/ca.csr \
        -subj "/C=CN/ST=Guangdong/L=Shenzhen/O=devops/OU=it/CN=127.0.0.1/emailAddress=nosbelm@qq.com" # 2. 生成请求文件
	@openssl x509 -req -in $(OUTPUT_DIR)/cert/ca.csr -signkey $(OUTPUT_DIR)/cert/ca.key -out $(OUTPUT_DIR)/cert/ca.crt # 3. 生成根证书
	@openssl genrsa -out $(OUTPUT_DIR)/cert/server.key 1024 # 4. 生成服务端私钥
	@openssl rsa -in $(OUTPUT_DIR)/cert/server.key -pubout -out $(OUTPUT_DIR)/cert/server.pem # 5. 生成服务端公钥
	@openssl req -new -key $(OUTPUT_DIR)/cert/server.key -out $(OUTPUT_DIR)/cert/server.csr \
        -subj "/C=CN/ST=Guangdong/L=Shenzhen/O=serverdevops/OU=serverit/CN=127.0.0.1/emailAddress=nosbelm@qq.com" # 6. 生成服务端向 CA 申请签名的 CSR
	@openssl x509 -req -CA $(OUTPUT_DIR)/cert/ca.crt -CAkey $(OUTPUT_DIR)/cert/ca.key \
        -CAcreateserial -in $(OUTPUT_DIR)/cert/server.csr -out $(OUTPUT_DIR)/cert/server.crt # 7. 生成服务端带有 CA 签名的证书

.PHONY: format
format: # 格式化go源码
	gofmt -s -w ./

.PHONY: swagger
swagger: # 启动swagger 在线文档
	@swagger serve -F=swagger --no-open --port 65543 $(ROOT_DIR)/api/openapi/miniblog/openapi.yaml

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: protoc
protoc: # 编译 protobuf 文件 https://juejin.cn/book/7176608782871429175/section/7179876228407492645#heading-7
	@echo "============> Generate protobuf files"
	@protoc                                            \
        --proto_path=$(APIROOT)                          \
        --proto_path=$(ROOT_DIR)/third_party             \
        --go_out=paths=source_relative:$(APIROOT)        \
        --go-grpc_out=paths=source_relative:$(APIROOT)   \
        $(shell find $(APIROOT) -name *.proto)


# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


