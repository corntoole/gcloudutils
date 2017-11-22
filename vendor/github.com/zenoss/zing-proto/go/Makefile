GLIDE := $(shell command -v glide 2> /dev/null)

glide.lock: glide.yaml
	@mkdir -p .glide
	@$(GLIDE) update --force --strip-vendor

vendor: glide.lock
	@mkdir -p .glide
	@$(GLIDE) install --force

.PHONY: clean
clean:
	rm -rf vendor
	rm -rf glide.lock