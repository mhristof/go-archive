MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
ifeq ($(word 1,$(subst ., ,$(MAKE_VERSION))),4)
.SHELLFLAGS := -eu -o pipefail -c
endif
.DEFAULT_GOAL := help
.ONESHELL:

GIT_REF := $(shell git describe --match="" --always --dirty=+)
GIT_TAG := $(shell git name-rev --tags --name-only $(GIT_REF))
PACKAGE := $(shell go list)

.PHONY: help
help:  ## Show this help
	@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ##/:/g' | column -t -s:

.git/hooks/pre-commit:  ## Install pre-commit checks
	pre-commit install

.PHONY: check
check: .git/hooks/pre-commit ## Run precommit checks
	pre-commit run --all

.PHONY: test
test:  ## Run go test
	go test -v ./...

bin/go-archive.darwin:  ## Build the application binary for current OS

bin/go-archive.%:  ## Build the application binary for target OS, for example bin/go-archive.linux
	GOOS=$* go build -o $@ -ldflags "-X $(PACKAGE)/version=$(GIT_TAG)+$(GIT_REF)" main.go

.PHONY: install
install: bin/go-archive.darwin ## Install the binary
	cp $< ~/bin/go-archive
