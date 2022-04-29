CURDIR = $(shell pwd)
GOBIN  = $(CURDIR)/bin

BINARY = dst
SRCDIR = $(CURDIR)
APPDIR = $(SRCDIR)

HOSTOS = $(shell go env GOHOSTOS)
HOSTARCH = $(shell go env GOHOSTARCH)

BUILDFLAGS = -trimpath -ldflags="-s -w"

default:
	unset GOPATH; \
	cd $(APPDIR); \
	GOOS=$(HOSTOS) GOARCH=$(HOSTARCH) go build $(BUILDFLAGS) -o $(GOBIN)/$(BINARY)

build:
	mkdir -p build
	for os in linux darwin windows; do \
		GOOS=$${os} GOARCH=amd64 go build $(BUILDFLAGS) -o build/$(BINARY)-$${os}; \
	done
	mv build/$(BINARY)-darwin build/$(BINARY)-osx
	mv build/$(BINARY)-windows build/$(BINARY)-win.exe

#build-linux:
#	for arch in 386 amd64 arm arm64 ppc64 ppc64le mips mipsle mips64 mips64le; do \
#		mkdir -p build/$${arch}; \
#		GOOS=linux GOARCH=$${arch} go build $(BUILDFLAGS) -o build/$${arch}/$(BINARY); \
#	done

.PHONY: default dist

