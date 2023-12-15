



include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk
include scripts/make-rules/gen.mk

.PHONY: init
# init .env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.15.1



## generate: Generate all necessary files, such as errors code files.
.PHONY: generate
generate:
	@$(MAKE) gen.run

.PHONY: protoc
protoc:
	@$(MAKE) gen.protoc


## ca: Generate CA files for all adminsrv components.
.PHONY: ca
ca:
	@$(MAKE) gen.ca


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
