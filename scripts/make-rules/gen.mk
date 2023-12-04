
.PHONY: gen.run
#gen.run: gen.errcode gen.docgo
gen.run: gen.clean gen.errcode


.PHONY: gen.errcode
gen.errcode: gen.errcode.code gen.errcode.doc

.PHONY: gen.errcode.code
gen.errcode.code: tools.verify.codegen
	@echo "===========> Generating voyage error code go source files"
	@codegen -type=int ${ROOT_DIR}/internal/pkg/code

.PHONY: gen.errcode.doc
gen.errcode.doc: tools.verify.codegen
	@echo "===========> Generating error code markdown documentation"
	@codegen -type=int -doc \
		-output ${ROOT_DIR}/docs/guide/zh-CN/api/error_code_generated.md ${ROOT_DIR}/internal/pkg/code



.PHONY: gen.clean
gen.clean:
	@rm -rf ./api/client/{clientset,informers,listers}
	@$(FIND) -type f -name '*_generated.go' -delete

.PHONY: gen.protoc
gen.protoc:
	@echo "=========> Generate protobuf files"
	@echo $(PROTO_DIR)
	@protoc 												\
		--proto_path=$(PROTO_DIR)							\
		--go_out=paths=source_relative:$(PROTO_GO_OUT_DIR) 		\
		--go-grpc_out=paths=source_relative:$(PROTO_GO_OUT_DIR) 	\
		--grpc-gateway_out=$(PROTO_GO_OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		$(shell find $(PROTO_DIR) -name *.proto)
