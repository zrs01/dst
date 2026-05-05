HOSTOS=$(shell go env GOHOSTOS)
HOSTARCH=$(shell go env GOHOSTARCH)

EXECUTABLE=dst
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
ARM=$(EXECUTABLE)_linux_arm
VERSION=$(shell git describe --tags --always --long --dirty)
BUILDFLAGS=-v -trimpath -ldflags="-s -w -X main.version=$(VERSION)"
CGO_ENABLED=0

# ---- default output filename -----------------------------------------
# If we are on Windows, add .exe; otherwise leave it without an extension
DEFAULT_BIN := bin/$(EXECUTABLE)
ifeq ($(HOSTOS),windows)
DEFAULT_BIN := $(DEFAULT_BIN).exe
endif
# ---------------------------------------------------------------------

default:
	GOOS=$(HOSTOS) GOARCH=$(HOSTARCH) CGO_ENABLED=$(CGO_ENABLED) go build -ldflags="-X main.version=$(VERSION)" -o $(DEFAULT_BIN)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

arm: $(ARM) ## Build for ARM 32-bit

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BUILDFLAGS) -o build/$(WINDOWS)

$(LINUX):
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BUILDFLAGS) -o build/$(LINUX)

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BUILDFLAGS) -o build/$(DARWIN)

$(ARM):
	env GOOS=linux GOARCH=arm CGO_ENABLED=$(CGO_ENABLED) go build $(BUILDFLAGS) -o build/$(ARM)

build: windows linux darwin arm ## Build binaries
	@echo version: $(VERSION)

# all: test build ## Build and run tests

# test: ## Run unit tests
# 	./scripts/test_unit.sh

clean: ## Remove previous build
	rm -f build/$(WINDOWS) build/$(LINUX) build/$(DARWIN) $(DEFAULT_BIN)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: default windows linux darwin arm clean help
