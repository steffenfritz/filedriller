ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=friller
VERSION=1.0.0-BETA-4
BUILD=`git rev-parse --short HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -w -s" 

all: clean build_all

build:
	go build ${LDFLAGS} -o ${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(LDFLAGS) -v -o $(BINARY)-$(GOOS))))


clean:
	rm -f '${BINARY}-linux'
	rm -f '${BINARY}-darwin'
	rm -f '${BINARY}-windows'

.PHONY: check clean install build_all all
