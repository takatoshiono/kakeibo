SHELL := /usr/bin/env bash -o pipefail

PROJECT := kakeibo

BUF_VERSION := 0.38.0
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)
CACHE_BASE := $(HOME)/.cache/$(PROJECT)
CACHE := $(CACHE_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
CACHE_BIN := $(CACHE)/bin

# Update the $PATH so we can use buf directly
export PATH := $(abspath $(CACHE_BIN)):$(PATH)

# install-buf is used in CI.
# ref: https://github.com/bufbuild/buf-example/blob/master/Makefile

.PHONY: install-buf
install-buf:
	@rm -f $(CACHE_BIN)/buf
	@mkdir -p $(CACHE_BIN)
	curl -sSL "https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/buf-$(UNAME_OS)-$(UNAME_ARCH)" -o "$(CACHE_BIN)/buf"
	chmod +x "$(CACHE_BIN)/buf"

# proto

.PHONY: proto-gen-go
proto:
	buf generate

# lint

.PHONY: lint
lint:
	buf lint
	buf breaking --against '.git#branch=master'

# clean deletes the cache for all platforms.

.PHONY: clean
clean:
	rm -rf $(CACHE_BASE)
