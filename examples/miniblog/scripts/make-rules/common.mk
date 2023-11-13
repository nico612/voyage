# ==============================================================================
# 定义全局 Makefile 变量方便后面引用
SHELL := /bin/bash

# ./ 根目录
# MAKEFILE_LIST 变量是 Makefile 的内置变量，表示：make 所需要处理的 makefile 文件列表，当前 makefile 的文件名总是位于列表的最后，文件名之间以空格进行分隔；
# 函数 $(lastword <text>) 取字符串 <text> 中的最后一个单词，并返回字符串 <text> 的最后一个单词；
# 函数 $(dir <names...>) 从文件名序列 <names> 中取出目录部分。目录部分是指最后一个反斜杠（/）之前的部分。如果没有反斜杠，那么返回 ./；
# 找出当前make文件所处的目录(绝对路径)
COMMON_SELF_DIR:=$(dir $(lastword $(MAKEFILE_LIST)))

# pwd -P 获取当前工作目录的物理路径（解析所有符号链接后的真实路径）
# 获取项目根目录真实路径。
# 获取当前执行的makefile文件所在路径： MAKE_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../../ && pwd -P))

# 构建产物，临时文件存放目录
OUTPUT_DIR := $(ROOT_DIR)/_output

# 定义包名
ROOT_PACKAGE=github.com/nico612/go-project

# Protobuf 文件存放路径
APIROOT=$(ROOT_DIR)/pkg/proto

# 创建TMP_DIR
ifeq ($(origin TMP_DIR), undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

# ==============================================================================
# 定义版本相关变量

## VERSION:=$(shell git describe --tags --always)
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

# 编译的操作系统可以是 linux/windows/darwin
## darwin_amd64: macOS 操作系统的 Intel 架构 64 位处理器（x86_64）。
## darwin_arm64: macOS 操作系统的 Apple Silicon M1 芯片的 ARM 架构。
## windows_amd64: Windows 操作系统的 Intel 架构 64 位处理器（x86_64）。
## linux_amd64 ：Linux 操作系统的 Intel 架构 64 位处理器（x86_64）。
## linux_arm64 ：Linux 操作系统的 ARM 架构 64 位处理器。
PLATFORMS ?= darwin_amd64 darwin_arm64 windows_amd64 linux_amd64 linux_arm64


# 设置一个指定的操作系统
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
	 	# 操作系统
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		# 系统架构
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# 构建镜像时，使用 linux 作为默认的 OS
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# Makefile 设置
ifndef V
MAKEFLAGS += --no-print-directory
endif

# Linux 命令设置
FIND := find . ! -path './third_party/*' ! -path './vendor/*'
XARGS := xargs --no-run-if-empty
