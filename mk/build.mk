##@ Building

SRC_GO := $(shell find . -name \*.go -print)
SRC_DATA := $(shell find internal/tsrv/snip -print)

PACKAGE_BASE = github.com/hardcaporg/hardcap/internal
LDFLAGS = "-X $(PACKAGE_BASE)/version.BuildCommit=$(shell git rev-parse --short HEAD) -X $(PACKAGE_BASE)/version.BuildTime=$(shell date +'%Y-%m-%d_%T')"

build: build/agent ## Build all binaries

prep:
	-mkdir build/

build/agent: prep $(SRC_GO) $(SRC_DATA) ## Build agent
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o build/agent ./cmd/hardcap-agent

.PHONY: strip
strip: build ## Strip debug information
	strip build/*

.PHONY: clean
clean: ## Clean build artifacts
	-rm -rf build/