
.PHONY: ca.gen.%
ca.gen.%:
	$(eval CA := $(word 1,$(subst ., ,$*)))
	@echo "===========> Generating CA files for $(CA)"
	@${ROOT_DIR}/scripts/gencerts.sh generate-iam-cert $(OUTPUT_DIR)/cert $(CA)


.PHONY: ca.gen
ca.gen: $(addprefix ca.gen., $(CERTIFICATES))
