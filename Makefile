SHELL := /usr/bin/env bash -o pipefail

PROJECT := kakeibo

# This controls the remote HTTPS git location to compare against for breaking changes in CI.
#
# Most CI providers only clone the branch under test and to a certain depth, so when
# running buf breaking in CI, it is generally preferable to compare against
# the remote repository directly.
#
# Basic authentication is available, see https://buf.build/docs/inputs#https for more details.
HTTPS_GIT := https://github.com/takatoshiono/kakeibo.git

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

# proto-stubs

.PHONY: proto-stubs
proto-stubs:
	buf generate

# proto-lint

.PHONY: proto-lint
proto-lint:
	buf lint

# proto-breaking is for local
# This does breaking change detection against our local git repository.

.PHONE: proto-breaking
proto-breaking:
	buf breaking --against '.git#branch=main'

# proto-breaking-ci is for CI
# This does breaking change detection against our remote HTTPS git repository.

.PHONE: proto-breaking-ci
proto-breaking-ci:
	buf breaking --against "$(HTTPS_GIT)#branch=main"

.PHONY: test-backend
test-backend:
	go test -race -v -coverprofile ./backend/coverage.out ./backend/...

.PHONY: build
build-backend:
	mkdir -p build
	go build -o ./backend/build ./backend/cmd/mf

# clean deletes the cache for all platforms.

.PHONY: clean
clean:
	rm -rf $(CACHE_BASE)
	rm -rf backend/build
