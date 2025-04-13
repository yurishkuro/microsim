# Copyright (c) 2023 The Jaeger Authors.
# SPDX-License-Identifier: Apache-2.0

SHELL := /bin/bash

# These DOCKER_xxx vars are used when building Docker images.
DOCKER_NAMESPACE?=yurishkuro
DOCKER_TAG?=latest

# SRC_ROOT is the top of the source tree.
SRC_ROOT := $(shell git rev-parse --show-toplevel)

GO=go
GOFMT=gofmt

# All .go files that are not auto-generated and should be auto-formatted and linted.
ALL_SRC = $(shell find . -name '*.go' \
				   -not -name '_*' \
				   -not -name '.*' \
				   -not -name 'mocks*' \
				   -not -name '*.pb.go' \
				   -not -path './vendor/*' \
				   -not -path './internal/tools/*' \
				   -not -path '*/mocks/*' \
				   -not -path '*/*-gen/*' \
				   -type f | \
				sort)

# import other Makefiles after the variables are defined
include Makefile.Tools.mk

.PHONY: fmt
fmt:
	@echo Running gofmt on ALL_SRC ...
	@$(GOFMT) -e -s -l -w $(ALL_SRC)
	@echo Running gofumpt on ALL_SRC ...
	@$(GOFUMPT) -e -l -w $(ALL_SRC)

.PHONY: lint
lint: $(LINT)
	$(LINT) -v run
