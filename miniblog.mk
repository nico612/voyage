GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION:=$(shell git describe --tags --always)

# ./ 根目录
# MAKEFILE_LIST 变量是 Makefile 的内置变量，表示：make 所需要处理的 makefile 文件列表，当前 makefile 的文件名总是位于列表的最后，文件名之间以空格进行分隔；
# 函数 $(lastword <text>) 取字符串 <text> 中的最后一个单词，并返回字符串 <text> 的最后一个单词；
# 函数 $(dir <names...>) 从文件名序列 <names> 中取出目录部分。目录部分是指最后一个反斜杠（/）之前的部分。如果没有反斜杠，那么返回 ./；
COMMON_SELF_DIR:=$(dir $(lastword $(MAKEFILE_LIST)))
# 获取项目根目录绝对路径。
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# bin 目录
BIN_DIR := $(ROOT_DIR)/bin/miniblog

# 构建产物，临时文件存放目录
OUTPUT_DIR := $(ROOT_DIR)/_output


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
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BIN_DIR)/miniblog $(ROOT_DIR)/cmd/miniblog/main.go

# .PHONY: run
# run:

.PHONY: generate
# generate
generate: tidy
	go generate ./...

.PHONY: clean # 实现幂等删除。
clean:
	-rm -vrf bin
	-rm -vrf api/doc
	-rm -vrf api/gen

.PHONY: format
format: # 格式化go源码
	gofmt -s -w ./

.PHONY: swagger
swagger: # 启动swagger 在线文档
	@swagger serve -F=swagger --no-open --port 65543 $(ROOT_DIR)/api/openapi/miniblog/openapi.yaml

.PHONY: tidy
tidy:
	go mod tidy
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


